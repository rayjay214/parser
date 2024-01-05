package ipc

import (
    _ "fmt"
)

type Body_0x8001 struct {
    Seq    uint16
    MsgId  MsgID
    Result byte
}

func (entity *Body_0x8001) MsgID() MsgID {
    return Msg_0x8001
}

func (entity *Body_0x8001) Encode() ([]byte, error) {
    writer := NewWriter()

    writer.WriteUint16(entity.Seq)

    writer.WriteUint16(uint16(entity.MsgId))

    writer.WriteByte(entity.Result)

    return writer.Bytes(), nil
}

func (entity *Body_0x8001) Decode(data []byte) (int, error) {
    reader := NewReader(data)

    var err error
    entity.Seq, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    var msgId uint16
    msgId, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }
    entity.MsgId = MsgID(msgId)

    entity.Result, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
