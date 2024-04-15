package ipc

import (
	_ "fmt"
)

type Body_0x8003 struct {
	AppId string
}

func (entity *Body_0x8003) MsgID() MsgID {
	return Msg_0x8003
}

func (entity *Body_0x8003) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.AppId, 20)

	return writer.Bytes(), nil
}

func (entity *Body_0x8003) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.AppId, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
