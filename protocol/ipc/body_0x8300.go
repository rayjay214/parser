package ipc

import (
    _ "fmt"
)

type Body_0x8300 struct {
    Encoding byte
    Len      uint16
    Content  string
}

func (entity *Body_0x8300) MsgID() MsgID {
    return Msg_0x8300
}

func (entity *Body_0x8300) Encode() ([]byte, error) {
    writer := NewWriter()

    writer.WriteByte(entity.Encoding)

    writer.WriteUint16(entity.Len)

    writer.WriteString(entity.Content, int(entity.Len))

    return writer.Bytes(), nil
}

func (entity *Body_0x8300) Decode(data []byte) (int, error) {
    reader := NewReader(data)

    var err error
    entity.Encoding, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Len, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    entity.Content, err = reader.ReadString(int(entity.Len))
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
