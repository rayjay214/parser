package ipc

import (
	_ "fmt"
)

type Body_0x8205 struct {
}

func (entity *Body_0x8205) MsgID() MsgID {
	return Msg_0x8205
}

func (entity *Body_0x8205) Encode() ([]byte, error) {
	writer := NewWriter()

	return writer.Bytes(), nil
}

func (entity *Body_0x8205) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	return len(data) - reader.Len(), nil
}
