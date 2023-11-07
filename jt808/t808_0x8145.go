package jt808

import (
    "github.com/rayjay214/parser/common"
)

// 终端应答
type T808_0x8145 struct {
}

func (entity *T808_0x8145) MsgID() MsgID {
    return MsgT808_0x8145
}

func (entity *T808_0x8145) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *T808_0x8145) Decode(data []byte) (int, error) {
    return 0, nil
}
