package jt808

import (
    "github.com/rayjay214/parser/common"
    "encoding/json"
)

// 请求同步时间
type T808_0x0105 struct {
    Type byte
}

func (entity T808_0x0105) MarshalJSON() ([]byte, error) {
    type Alias T808_0x0105

    type NewT808_0x0105 struct {
        Alias
        Type string
    }

    s := NewT808_0x0105{
        Alias: Alias(entity),
    }

    switch entity.Type {
    case 0:
        s.Type = "休眠（网络连接正常，关闭GPS）"
    case 1:
        s.Type = "休眠（网络连接正常，打开GPS）"
    case 2:
        s.Type = "休眠（断网且关GPS）"
    default:
        s.Type = "其他"
    }

    return json.Marshal(s)
}

func (entity *T808_0x0105) MsgID() MsgID {
    return MsgT808_0x1006
}

func (entity *T808_0x0105) Encode() ([]byte, error) {
    return nil, nil
}

func (entity *T808_0x0105) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    entity.Type, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    return 0, nil
}
