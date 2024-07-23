package ipc

import (
	_ "fmt"
)

type Body_0x0011 struct {
	DeviceType string
}

func (entity *Body_0x0011) MsgID() MsgID {
	return Msg_0x0011
}

func (entity *Body_0x0011) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.DeviceType, 20)

	return writer.Bytes(), nil
}

func (entity *Body_0x0011) Decode(data []byte) (int, error) {
	reader := NewReader(data)
	var err error

	entity.DeviceType, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
