package ipc

import (
    _ "fmt"
)

type Body_0x0004 struct {
    Uid string
}

func (entity *Body_0x0004) MsgID() MsgID {
    return Msg_0x0004
}

func (entity *Body_0x0004) Encode() ([]byte, error) {
    writer := NewWriter()

    writer.WriteString(entity.Uid, 20)

    return writer.Bytes(), nil
}

func (entity *Body_0x0004) Decode(data []byte) (int, error) {
    reader := NewReader(data)

    var err error
    entity.Uid, err = reader.ReadString(20)
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
