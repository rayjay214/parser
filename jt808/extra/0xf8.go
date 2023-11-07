package extra

import (
    "parser/jt808/errors"
    "encoding/binary"
)

type Extra_0xf8 struct {
    serialized []byte
    value      uint16
}

func (Extra_0xf8) ID() byte {
    return byte(TypeExtra_0xf8)
}

func (extra Extra_0xf8) Data() []byte {
    return extra.serialized
}

func (extra Extra_0xf8) Value() interface{} {
    return extra.value
}

func (extra Extra_0xf8) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "计步步数"
    m["value"] = extra.Value()

    return m
}

func (extra *Extra_0xf8) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidExtraLength
    }
    extra.serialized = data
    extra.value = binary.BigEndian.Uint16(data)
    return 2, nil
}
