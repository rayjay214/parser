package extra

import (
	"github.com/rayjay214/parser/protocol/jt808/errors"
)

type Extra_0xf9 struct {
	serialized []byte
	value      byte
}

func (Extra_0xf9) ID() byte {
	return byte(TypeExtra_0xf9)
}

func (extra Extra_0xf9) Data() []byte {
	return extra.serialized
}

func (extra Extra_0xf9) Value() interface{} {
	return extra.value
}

func (extra Extra_0xf9) ToPrint() interface{} {
	m := map[string]interface{}{}
	m["desc"] = "是否上报假点"

	if extra.value == 0 {
		m["value"] = "不上报"
	} else {
		m["value"] = "上报"
	}

	return m
}

func (extra *Extra_0xf9) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.serialized = data
	extra.value = data[0]
	return 1, nil
}
