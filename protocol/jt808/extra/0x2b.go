package extra

import (
	"encoding/binary"
	"github.com/rayjay214/parser/protocol/jt808/errors"
)

// 基站信息
type Extra_0x2b struct {
	serialized []byte
	value      Extra_0x2b_Value
}

type Extra_0x2b_Value struct {
	Voltage1 uint16 `json:"-"`
	Voltage2 uint16 `json:"-"`
}

func (Extra_0x2b) ID() byte {
	return byte(TypeExtra_0x2b)
}

func (extra Extra_0x2b) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x2b) Value() interface{} {
	return extra.value
}

func (extra Extra_0x2b) ToPrint() interface{} {
	m := map[string]interface{}{}
	m["desc"] = "电压"
	m["value"] = extra.Value()

	return m
}

func (extra *Extra_0x2b) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidExtraLength
	}
	voltage1 := binary.BigEndian.Uint16(data)
	voltage2 := binary.BigEndian.Uint16(data[2:])

	extra.value.Voltage1 = voltage1
	extra.value.Voltage2 = voltage2

	extra.serialized = data
	return len(data), nil
}
