package extra

import (
	"github.com/rayjay214/parser/protocol/jt808/errors"
	"encoding/json"
	"fmt"
)

// 基站信息
type Extra_0xe4 struct {
	serialized []byte
	value      Extra_0xe4_Value
}

type Extra_0xe4_Value struct {
	Status byte `json:"-"`
	Power  byte `json:"-"`
}

func (entity Extra_0xe4_Value) MarshalJSON() ([]byte, error) {
	type Alias Extra_0xe4_Value

	type NewExtra_0xe4_Value struct {
		Alias
		Power  string `json:"电量"`
		Status string `json:"状态"`
	}

	s := NewExtra_0xe4_Value{
		Alias: Alias(entity),
	}

	s.Power = fmt.Sprintf("%d", entity.Power) + "%"

	if entity.Status == 0 {
		s.Status = "接通外接电源"
	} else {
		s.Status = "断开外接电源"
	}

	return json.Marshal(s)
}

func (Extra_0xe4) ID() byte {
	return byte(TypeExtra_0xe4)
}

func (extra Extra_0xe4) Data() []byte {
	return extra.serialized
}

func (extra Extra_0xe4) Value() interface{} {
	return extra.value
}

func (extra Extra_0xe4) ToPrint() interface{} {
	m := map[string]interface{}{}
	m["desc"] = "状态，电量"
	m["value"] = extra.Value()

	return m
}

func (extra *Extra_0xe4) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidExtraLength
	}
	extra.value.Status = data[0]
	extra.value.Power = data[1]

	extra.serialized = data
	return len(data), nil
}
