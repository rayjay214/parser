package ipc

import (
    _ "fmt"
)

type Body_0x0002 struct {
    Info string
}

func (entity *Body_0x0002) MsgID() MsgID {
    return Msg_0x0002
}

func (entity *Body_0x0002) Encode() ([]byte, error) {
    writer := NewWriter()

    writer.WriteString(entity.Info, 30)

    return writer.Bytes(), nil
}

func (entity *Body_0x0002) Decode(data []byte) (int, error) {
    reader := NewReader(data)

    var err error
    entity.Info, err = reader.ReadString(30)
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
