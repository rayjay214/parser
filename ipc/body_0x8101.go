package ipc

import (
	_ "fmt"
)

type Body_0x8101 struct {
	AppId string
	Ip    string
	Port  uint16
}

func (entity *Body_0x8101) MsgID() MsgID {
	return Msg_0x8101
}

func (entity *Body_0x8101) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.AppId, 20)

	writer.WriteString(entity.Ip, 20)

	writer.WriteUint16(entity.Port)

	return writer.Bytes(), nil
}

func (entity *Body_0x8101) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.AppId, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.Ip, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.Port, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
