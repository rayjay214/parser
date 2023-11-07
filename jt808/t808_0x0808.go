package jt808

import (
    "github.com/rayjay214/parser/common"
)

// 终端应答
type T808_0x0808 struct {
}

func (entity *T808_0x0808) MsgID() MsgID {
    return MsgT808_0x0808
}

func (entity *T808_0x0808) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *T808_0x0808) Decode(data []byte) (int, error) {
    return 0, nil
}
