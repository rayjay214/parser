package ipc

import (
    _ "fmt"
)

type Body_0x8004 struct {
    Seq  uint16
    Ip   string
    Port uint16
}

func (entity *Body_0x8004) MsgID() MsgID {
    return Msg_0x8004
}

func (entity *Body_0x8004) Encode() ([]byte, error) {
    writer := NewWriter()

    writer.WriteUint16(entity.Seq)

    writer.WriteString(entity.Ip, 20)

    writer.WriteUint16(entity.Port)

    return writer.Bytes(), nil
}

func (entity *Body_0x8004) Decode(data []byte) (int, error) {
    reader := NewReader(data)

    var err error
    entity.Seq, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    entity.Ip, err = reader.ReadString(20)
    if err != nil {
        return 0, err
    }

    entity.Port, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
