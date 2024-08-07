package ipc

import (
    _ "fmt"
)

type Body_0x8004 struct {
    AppId     string
    EventType uint16
}

func (entity *Body_0x8004) MsgID() MsgID {
    return Msg_0x8004
}

func (entity *Body_0x8004) Encode() ([]byte, error) {
    writer := NewWriter()

    writer.WriteString(entity.AppId, 20)

    writer.WriteUint16(entity.EventType)

    return writer.Bytes(), nil
}

func (entity *Body_0x8004) Decode(data []byte) (int, error) {
    reader := NewReader(data)

    var err error

    entity.AppId, err = reader.ReadString(20)
    if err != nil {
        return 0, err
    }

    entity.EventType, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
