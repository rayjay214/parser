package extra

import (
    "parser/jt808/errors"
)

// 基站信息
type Extra_0xe6 struct {
    serialized []byte
    value      byte
}

/*
func NewExtra_0xe1(val byte) *Extra_0xe1 {
	extra := Extra_0xe1{
		value: val,
	}
	extra.serialized = []byte{val}
	return &extra
}
*/

func (Extra_0xe6) ID() byte {
    return byte(TypeExtra_0xe6)
}

func (extra Extra_0xe6) Data() []byte {
    return extra.serialized
}

func (extra Extra_0xe6) Value() interface{} {
    return extra.value
}

func (extra Extra_0xe6) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "是否发短信"

    if extra.value == 0 {
        m["value"] = "该定位不用发短信"
    } else {
        m["value"] = "该定位要发短信使用"
    }

    return m
}

func (extra *Extra_0xe6) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidExtraLength
    }
    extra.value = data[0]

    extra.serialized = data
    return len(data), nil
}
