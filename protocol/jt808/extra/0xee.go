package extra

import (
	"github.com/rayjay214/parser/protocol/jt808/errors"
	"encoding/hex"
)

type Extra_0xee struct {
	serialized []byte
	//todo
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

func (Extra_0xee) ID() byte {
	return byte(TypeExtra_0xee)
}

func (extra Extra_0xee) Data() []byte {
	return extra.serialized
}

func (extra Extra_0xee) Value() interface{} {
	return nil
}

func (extra Extra_0xee) ToPrint() interface{} {
	m := map[string]interface{}{}
	m["desc"] = "4G基站信息"

	m["value"] = hex.EncodeToString(extra.serialized)

	return m
}

func (extra *Extra_0xee) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidExtraLength
	}

	extra.serialized = data
	return len(data), nil
}
