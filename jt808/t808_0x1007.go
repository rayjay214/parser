package jt808

import (
    "github.com/rayjay214/parser/common"
    "encoding/json"
    "fmt"
)

// 请求同步时间
type T808_0x1007 struct {
    Timezone int `json:"-"`
    Reserve1 byte
    Reserve2 uint16
}

func (entity *T808_0x1007) MsgID() MsgID {
    return MsgT808_0x1007
}

func (entity *T808_0x1007) Encode() ([]byte, error) {
    return nil, nil
}

func (entity *T808_0x1007) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    var timezone byte
    timezone, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }
    entity.Timezone = int(timezone) / 8
    if common.GetBit(int(timezone), 0) != 0 {
        entity.Timezone = entity.Timezone * -1
    } else {
        entity.Timezone = entity.Timezone
    }

    entity.Reserve1, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Reserve2, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    return 0, nil
}

func (entity T808_0x1007) MarshalJSON() ([]byte, error) {
    type Alias T808_0x1007

    type New1007 struct {
        Alias
        Timezone string `json:"时区"`
    }

    s := New1007{
        Alias: Alias(entity),
    }

    if entity.Timezone > 0 {
        s.Timezone = fmt.Sprintf("东%d区", entity.Timezone)
    } else if entity.Timezone < 0 {
        s.Timezone = fmt.Sprintf("西%d区", entity.Timezone*-1)
    } else {
        s.Timezone = "零时区"
    }

    return json.Marshal(s)
}
