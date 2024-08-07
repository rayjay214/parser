package extra

import (
	"github.com/rayjay214/parser/protocol/jt808/errors"
	"encoding/binary"
	"fmt"
)

type Extra_0xf1 struct {
	serialized []byte
	value      uint16
}

func (Extra_0xf1) ID() byte {
	return byte(TypeExtra_0xf1)
}

func (extra Extra_0xf1) Data() []byte {
	return extra.serialized
}

func (extra Extra_0xf1) Value() interface{} {
	return extra.value
}

func (extra Extra_0xf1) ToPrint() interface{} {
	m := map[string]interface{}{}
	m["desc"] = "当前电流，单位mA"
	m["value"] = fmt.Sprintf("%dmA", extra.Value())

	return m
}

func (extra *Extra_0xf1) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.serialized = data
	extra.value = binary.BigEndian.Uint16(data)
	return 2, nil
}
