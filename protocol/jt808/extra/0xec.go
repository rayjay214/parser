package extra

import (
	"github.com/rayjay214/parser/protocol/jt808/errors"
	"encoding/hex"
)

type Extra_0xec struct {
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

func (Extra_0xec) ID() byte {
	return byte(TypeExtra_0xec)
}

func (extra Extra_0xec) Data() []byte {
	return extra.serialized
}

func (extra Extra_0xec) Value() interface{} {
	return nil
}

func (extra Extra_0xec) ToPrint() interface{} {
	m := map[string]interface{}{}
	m["desc"] = "WIFI信息"
	m["value"] = hex.EncodeToString(extra.serialized)

	return m
}

func (extra *Extra_0xec) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidExtraLength
	}

	extra.serialized = data
	return len(data), nil
}
