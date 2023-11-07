package extra

import (
    "github.com/rayjay214/parser/jt808/errors"
)

type Extra_0xf3 struct {
    serialized []byte
    value      byte
}

/*
func NewExtra_0xf3(val byte) *Extra_0xf3 {
	extra := Extra_0xf3{
		value: val,
	}
	extra.serialized = []byte{val}
	return &extra
}
*/

func (Extra_0xf3) ID() byte {
    return byte(TypeExtra_0xf3)
}

func (extra Extra_0xf3) Data() []byte {
    return extra.serialized
}

func (extra Extra_0xf3) Value() interface{} {
    return extra.value
}

func (extra Extra_0xf3) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "北斗卫星个数"

    if extra.Value() != nil {
        m["value"] = extra.Value()
    } else {
        m["value"] = extra.serialized
    }

    return m
}

func (extra *Extra_0xf3) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidExtraLength
    }
    extra.serialized = data
    extra.value = data[0]
    return 1, nil
}
