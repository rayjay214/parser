package extra

import (
    "encoding/binary"
    "fmt"
    "github.com/rayjay214/parser/jt808/errors"
)

// 油量
type Extra_0x02 struct {
    serialized []byte
    value      uint16
}

func NewExtra_0x02(val uint16) *Extra_0x02 {
    extra := Extra_0x02{
        value: val,
    }

    var temp [2]byte
    binary.BigEndian.PutUint16(temp[:2], val)
    extra.serialized = temp[:2]
    return &extra
}

func (Extra_0x02) ID() byte {
    return byte(TypeExtra_0x02)
}

func (extra Extra_0x02) Data() []byte {
    return extra.serialized
}

func (extra Extra_0x02) Value() interface{} {
    return extra.value
}

func (extra Extra_0x02) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "油量"
    m["value"] = fmt.Sprintf("%.2fL", float32(extra.value)/10)

    return m
}

func (extra *Extra_0x02) Decode(data []byte) (int, error) {
    if len(data) < 2 {
        return 0, errors.ErrInvalidExtraLength
    }
    extra.serialized = data
    extra.value = binary.BigEndian.Uint16(data)
    return 2, nil
}
