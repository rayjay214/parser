package ipc

import (
    _ "fmt"
)

type Body_0x0003 struct {
	Ip          string
	Port        uint16
	PrivateIp   string
	PrivatePort uint16
}

func (entity *Body_0x0003) MsgID() MsgID {
    return Msg_0x0003
}

func (entity *Body_0x0003) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Ip, 20)

	writer.WriteUint16(entity.Port)

	writer.WriteString(entity.PrivateIp, 20)

	writer.WriteUint16(entity.PrivatePort)

	return writer.Bytes(), nil
}

func (entity *Body_0x0003) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Ip, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.Port, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.PrivateIp, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.PrivatePort, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
