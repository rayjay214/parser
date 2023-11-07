package extra

import (
    "github.com/rayjay214/parser/jt808/errors"
    "encoding/hex"
)

// 基站信息
type Extra_0xe1 struct {
    serialized []byte
    //todo
}

func NewExtra_0xe1(hexStr string) *Extra_0xe1 {
    extra := Extra_0xe1{}

    data, _ := hex.DecodeString(hexStr)
    extra.serialized = data
    return &extra
}

func (Extra_0xe1) ID() byte {
    return byte(TypeExtra_0xe1)
}

func (extra Extra_0xe1) Data() []byte {
    return extra.serialized
}

func (extra Extra_0xe1) Value() interface{} {
    //todo
    return nil
}

func (extra Extra_0xe1) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "基站信息或电量"
    m["value"] = hex.EncodeToString(extra.serialized)

    return m
}

func (extra *Extra_0xe1) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidExtraLength
    }
    extra.serialized = data
    return len(data), nil
}
