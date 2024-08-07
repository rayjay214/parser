package extra

import (
	_ "encoding/binary"
	"github.com/rayjay214/parser/protocol/common"
	_ "github.com/rayjay214/parser/protocol/jt808/errors"
	"strconv"
)

type Extra_0xf6 struct {
	serialized []byte
	value      Extra_0xf6_Value
}

type Extra_0xf6_Value struct {
	Imei uint64
}

/*
func NewExtra_0xf6(val byte) *Extra_0xf6 {
	extra := Extra_0xf6{
		value: val,
	}
	extra.serialized = []byte{val}
	return &extra
}
*/

func (Extra_0xf6) ID() byte {
	return byte(TypeExtra_0xf6)
}

func (extra Extra_0xf6) Data() []byte {
	return extra.serialized
}

func (extra Extra_0xf6) Value() interface{} {
	return extra.value
}

func (extra Extra_0xf6) ToPrint() interface{} {
	m := map[string]interface{}{}
	m["desc"] = "IMEI上传"

	m["imei"] = extra.value.Imei

	return m
}

func (extra *Extra_0xf6) Decode(data []byte) (int, error) {

	extra.serialized = data

	reader := common.NewReader(data)
	temp, _ := reader.Read(8)
	extra.value.Imei, _ = strconv.ParseUint(common.BcdToString(temp), 10, 64)

	return 8, nil
}
