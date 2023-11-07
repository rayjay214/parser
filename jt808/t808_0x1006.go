package jt808

import (
    "github.com/rayjay214/parser/common"
)

// 请求同步时间
type T808_0x1006 struct {
    IsSleepCloseNet byte
    Reserve         uint16
}

func (entity *T808_0x1006) MsgID() MsgID {
    return MsgT808_0x1006
}

func (entity *T808_0x1006) Encode() ([]byte, error) {
    return nil, nil
}

func (entity *T808_0x1006) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    entity.IsSleepCloseNet, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Reserve, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    return 0, nil
}
