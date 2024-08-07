package extra

import (
	_ "github.com/rayjay214/parser/protocol/jt808/errors"
	"encoding/binary"
	"fmt"
	"encoding/json"
)

type Extra_0x51 struct {
	serialized []byte
	value      Extra_0x51_Value
}

type Extra_0x51_Value struct {
	Temperature1 uint16 `json:"-"`
	Temperature2 uint16 `json:"-"`
	Temperature3 uint16 `json:"-"`
	Temperature4 uint16 `json:"-"`
}

func (entity Extra_0x51_Value) MarshalJSON() ([]byte, error) {
	type Alias Extra_0x51_Value

	type NewExtra_0x51_Value struct {
		Alias
		Temperature1 string
		Temperature2 string
		Temperature3 string
		Temperature4 string
	}

	s := NewExtra_0x51_Value{
		Alias: Alias(entity),
	}

	if entity.Temperature1 > 0x8000 {
		s.Temperature1 = fmt.Sprintf("%.2f摄氏度", float32(entity.Temperature1-0x8000)*(-1)/10)
	} else {
		s.Temperature1 = fmt.Sprintf("%.2f摄氏度", float32(entity.Temperature1)/10)
	}
	if entity.Temperature2 > 0x8000 {
		s.Temperature2 = fmt.Sprintf("%.2f摄氏度", float32(entity.Temperature2-0x8000)*(-1)/10)
	} else {
		s.Temperature2 = fmt.Sprintf("%.2f摄氏度", float32(entity.Temperature2)/10)
	}
	if entity.Temperature3 > 0x8000 {
		s.Temperature3 = fmt.Sprintf("%.2f摄氏度", float32(entity.Temperature3-0x8000)*(-1)/10)
	} else {
		s.Temperature3 = fmt.Sprintf("%.2f摄氏度", float32(entity.Temperature3)/10)
	}
	if entity.Temperature4 > 0x8000 {
		s.Temperature4 = fmt.Sprintf("%.2f摄氏度", float32(entity.Temperature4-0x8000)*(-1)/10)
	} else {
		s.Temperature4 = fmt.Sprintf("%.2f摄氏度", float32(entity.Temperature4)/10)
	}

	return json.Marshal(s)
}

func (Extra_0x51) ID() byte {
	return byte(TypeExtra_0x51)
}

func (extra Extra_0x51) Data() []byte {
	return extra.serialized
}

func (extra Extra_0x51) Value() interface{} {
	return extra.value
}

func (extra Extra_0x51) ToPrint() interface{} {
	m := map[string]interface{}{}
	m["desc"] = "4路温度上传"
	m["value"] = extra.Value()

	return m
}

func (extra *Extra_0x51) Decode(data []byte) (int, error) {

	extra.serialized = data
	extra.value.Temperature1 = binary.BigEndian.Uint16(data)
	extra.value.Temperature2 = binary.BigEndian.Uint16(data[2:])
	extra.value.Temperature3 = binary.BigEndian.Uint16(data[4:])
	extra.value.Temperature4 = binary.BigEndian.Uint16(data[6:])

	return 8, nil
}
