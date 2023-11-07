package jt808

import (
    "parser/common"
)

// 终端应答
type T808_0x8155 struct {
}

func (entity *T808_0x8155) MsgID() MsgID {
    return MsgT808_0x8155
}

func (entity *T808_0x8155) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *T808_0x8155) Decode(data []byte) (int, error) {
    return 0, nil
}
