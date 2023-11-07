package kks

import (
    "encoding/hex"
    "encoding/json"
    "fmt"
    "parser/common"
    _ "strconv"
)

type Kks_0x94 struct {
    Proto    uint8
    SubProto uint8
    Content  []byte
}

func (entity *Kks_0x94) MsgID() MsgID {
    return Msg_0x94
}

func (entity *Kks_0x94) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *Kks_0x94) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.Proto, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.SubProto, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    contentLen := len(data) - 8
    entity.Content, err = reader.Read(contentLen)
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}

func (entity Kks_0x94) MarshalJSON() ([]byte, error) {
    type Alias Kks_0x94

    type New0x94 struct {
        Proto    string
        SubProto string
        Content  string
        Alias
    }

    s := New0x94{
        Alias: Alias(entity),
    }

    s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

    switch entity.SubProto {
    case 0x00:
        s.SubProto = "外电电压"
        s.Content = hex.EncodeToString(entity.Content)
    case 0x04:
        s.SubProto = "终端状态同步"
        s.Content = string(entity.Content)
    case 0x05:
        s.SubProto = "门状态"
        s.Content = string(entity.Content)
    case 0x08:
        s.SubProto = "自检参数"
        s.Content = string(entity.Content)
    case 0x09:
        s.SubProto = "定位卫星信息"
        s.Content = hex.EncodeToString(entity.Content)
    case 0x0A:
        s.SubProto = "ICCID信息"
        s.Content = string(entity.Content)
    }

    return json.Marshal(s)
}
