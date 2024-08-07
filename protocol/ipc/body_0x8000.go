package ipc

import (
	_ "fmt"
)

type Body_0x8000 struct {
	Seq       uint16
	RelayIp   string //转发服务IP
	RelayPort uint16 //转发服务端口
	RtpPort   uint16 //直推端口
}

func (entity *Body_0x8000) MsgID() MsgID {
	return Msg_0x8000
}

func (entity *Body_0x8000) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteUint16(entity.Seq)

	writer.WriteString(entity.RelayIp, 20)

	writer.WriteUint16(entity.RelayPort)

	writer.WriteUint16(entity.RtpPort)

	return writer.Bytes(), nil
}

func (entity *Body_0x8000) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Seq, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.RelayIp, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.RelayPort, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.RtpPort, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
