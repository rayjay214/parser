package extra

import (
    "github.com/rayjay214/parser/jt808/errors"
    "encoding/hex"
)

// wifi信息
type Extra_0x54 struct {
    serialized []byte
    //todo
}

func NewExtra_0x54(hexStr string) *Extra_0x54 {
    extra := Extra_0x54{}

    data, _ := hex.DecodeString(hexStr)
    extra.serialized = data
    return &extra
}

func (Extra_0x54) ID() byte {
    return byte(TypeExtra_0x54)
}

func (extra Extra_0x54) Data() []byte {
    return extra.serialized
}

func (extra Extra_0x54) Value() interface{} {
    //todo
    return nil
}

func (extra Extra_0x54) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "移文wifi信息"
    m["value"] = hex.EncodeToString(extra.serialized)

    return m
}

func (extra *Extra_0x54) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidExtraLength
    }
    extra.serialized = data
    return len(data), nil
}
