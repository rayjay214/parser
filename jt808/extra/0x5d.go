package extra

import (
    "encoding/hex"
    "parser/jt808/errors"
)

// 基站信息
type Extra_0x5d struct {
    serialized []byte
    //todo
}

func NewExtra_0x5d(hexStr string) *Extra_0x5d {
    extra := Extra_0x5d{}

    data, _ := hex.DecodeString(hexStr)
    extra.serialized = data
    return &extra
}

func (Extra_0x5d) ID() byte {
    return byte(TypeExtra_0x5d)
}

func (extra Extra_0x5d) Data() []byte {
    return extra.serialized
}

func (extra Extra_0x5d) Value() interface{} {
    //todo
    return nil
}

func (extra Extra_0x5d) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "移文基站信息"
    m["value"] = hex.EncodeToString(extra.serialized)

    return m
}

func (extra *Extra_0x5d) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidExtraLength
    }
    extra.serialized = data
    return len(data), nil
}
