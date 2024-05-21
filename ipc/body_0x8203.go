package ipc

import (
	_ "fmt"
)

type Body_0x8203 struct {
	Result uint8
}

func (entity *Body_0x8203) MsgID() MsgID {
	return Msg_0x8203
}

func (entity *Body_0x8203) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteByte(entity.Result)

	return writer.Bytes(), nil
}

func (entity *Body_0x8203) Decode(data []byte) (int, error) {
	reader := NewReader(data)
	var err error
	entity.Result, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
