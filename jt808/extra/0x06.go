package extra

import (
    "github.com/rayjay214/parser/jt808/errors"
)

// 基站信息
type Extra_0x06 struct {
    serialized []byte
    value      byte
}

func (Extra_0x06) ID() byte {
    return byte(TypeExtra_0x06)
}

func (extra Extra_0x06) Data() []byte {
    return extra.serialized
}

func (extra Extra_0x06) Value() interface{} {
    return extra.value
}

func (extra Extra_0x06) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "是否发短信"

    if extra.value == 0 {
        m["value"] = "该定位不用发短信"
    } else {
        m["value"] = "该定位要发短信使用"
    }

    return m
}

func (extra *Extra_0x06) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidExtraLength
    }
    extra.value = data[0]

    extra.serialized = data
    return len(data), nil
}
