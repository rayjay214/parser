package ipc

import (
	_ "fmt"
)

type Body_0x0021 struct {
	Info string
}

func (entity *Body_0x0021) MsgID() MsgID {
	return Msg_0x0021
}

func (entity *Body_0x0021) Encode() ([]byte, error) {
	writer := NewWriter()

	return writer.Bytes(), nil
}

func (entity *Body_0x0021) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	return len(data) - reader.Len(), nil
}
