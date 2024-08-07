package extra

import (
	"github.com/rayjay214/parser/protocol/jt808/errors"
)

type Extra_0xf2 struct {
	serialized []byte
	value      byte
}

/*
func NewExtra_0xf2(val byte) *Extra_0xf2 {
	extra := Extra_0xf2{
		value: val,
	}
	extra.serialized = []byte{val}
	return &extra
}
*/

func (Extra_0xf2) ID() byte {
	return byte(TypeExtra_0xf2)
}

func (extra Extra_0xf2) Data() []byte {
	return extra.serialized
}

func (extra Extra_0xf2) Value() interface{} {
	return extra.value
}

func (extra Extra_0xf2) ToPrint() interface{} {
	m := map[string]interface{}{}
	m["desc"] = "GPS卫星个数"

	if extra.Value() != nil {
		m["value"] = extra.Value()
	} else {
		m["value"] = extra.serialized
	}

	return m
}

func (extra *Extra_0xf2) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.serialized = data
	extra.value = data[0]
	return 1, nil
}
