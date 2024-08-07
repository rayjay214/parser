package extra

import (
	"github.com/rayjay214/parser/protocol/jt808/errors"
	"encoding/binary"
	"fmt"
)

type Extra_0xf0 struct {
	serialized []byte
	value      uint16
}

func (Extra_0xf0) ID() byte {
	return byte(TypeExtra_0xf0)
}

func (extra Extra_0xf0) Data() []byte {
	return extra.serialized
}

func (extra Extra_0xf0) Value() interface{} {
	return extra.value
}

func (extra Extra_0xf0) ToPrint() interface{} {
	m := map[string]interface{}{}
	m["desc"] = "外电电压"
	m["value"] = fmt.Sprintf("%.2fV", float32(extra.value)/100)

	return m
}

func (extra *Extra_0xf0) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.serialized = data
	extra.value = binary.BigEndian.Uint16(data)
	return 2, nil
}
