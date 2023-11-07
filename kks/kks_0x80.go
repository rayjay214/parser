package kks

import (
    "encoding/json"
    "fmt"
    "parser/common"
    _ "strconv"
)

type Kks_0x80 struct {
    Proto      uint8
    ContentLen uint8
    SysFlag    uint32
    Content    string
    Reserved   uint8
    Lang       uint8
}

func (entity *Kks_0x80) MsgID() MsgID {
    return Msg_0x80
}

func (entity *Kks_0x80) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *Kks_0x80) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.Proto, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.ContentLen, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.SysFlag, err = reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    entity.Content, err = reader.ReadString(int(entity.ContentLen - 4))
    if err != nil {
        return 0, err
    }

    entity.Reserved, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Lang, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}

func (entity Kks_0x80) MarshalJSON() ([]byte, error) {
    type Alias Kks_0x80

    type Kks_0x80 struct {
        Proto string
        Alias
    }

    s := Kks_0x80{
        Alias: Alias(entity),
    }

    s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

    return json.Marshal(s)
}
