package extra

import (
    "encoding/binary"
    "fmt"
    "github.com/rayjay214/parser/jt808/errors"
)

// 里程
type Extra_0x01 struct {
    serialized []byte
    value      uint32
}

func NewExtra_0x01(val uint32) *Extra_0x01 {
    extra := Extra_0x01{
        value: val,
    }

    var temp [4]byte
    binary.BigEndian.PutUint32(temp[:4], val)
    extra.serialized = temp[:4]
    return &extra
}

func (Extra_0x01) ID() byte {
    return byte(TypeExtra_0x01)
}

func (extra Extra_0x01) Data() []byte {
    return extra.serialized
}

func (extra Extra_0x01) Value() interface{} {
    return extra.value
}

func (extra Extra_0x01) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "里程"
    m["value"] = fmt.Sprintf("%.2fkm", float32(extra.value)/10)

    return m
}

func (extra *Extra_0x01) Decode(data []byte) (int, error) {
    if len(data) < 4 {
        return 0, errors.ErrInvalidExtraLength
    }
    extra.serialized = data
    extra.value = binary.BigEndian.Uint32(data)
    return 4, nil
}
