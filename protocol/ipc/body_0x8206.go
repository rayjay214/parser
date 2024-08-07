package ipc

import (
	_ "fmt"
)

type Body_0x8206 struct {
	Appid string
	Month uint32
}

func (entity *Body_0x8206) MsgID() MsgID {
	return Msg_0x8206
}

func (entity *Body_0x8206) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Appid, 20)

	writer.WriteUint32(entity.Month)

	return writer.Bytes(), nil
}

func (entity *Body_0x8206) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Appid, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.Month, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
