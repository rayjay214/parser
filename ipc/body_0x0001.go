package ipc

import (
    _ "fmt"
)

type Body_0x0001 struct {
}

func (entity *Body_0x0001) MsgID() MsgID {
    return Msg_0x0001
}

func (entity *Body_0x0001) Encode() ([]byte, error) {
    writer := NewWriter()

    return writer.Bytes(), nil
}

func (entity *Body_0x0001) Decode(data []byte) (int, error) {
    reader := NewReader(data)

    return len(data) - reader.Len(), nil
}
