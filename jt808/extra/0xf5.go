package extra

import (
    "parser/jt808/errors"
)

type Extra_0xf5 struct {
    serialized []byte
    value      byte
}

/*
func NewExtra_0xf5(val byte) *Extra_0xf5 {
	extra := Extra_0xf5{
		value: val,
	}
	extra.serialized = []byte{val}
	return &extra
}
*/

func (Extra_0xf5) ID() byte {
    return byte(TypeExtra_0xf5)
}

func (extra Extra_0xf5) Data() []byte {
    return extra.serialized
}

func (extra Extra_0xf5) Value() interface{} {
    return extra.value
}

func (extra Extra_0xf5) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "信号类型"

    if extra.value == 0 {
        m["value"] = "2G信号上传"
    } else {
        m["value"] = "4G信号上传"
    }

    return m
}

func (extra *Extra_0xf5) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidExtraLength
    }
    extra.serialized = data
    extra.value = data[0]
    return 1, nil
}
