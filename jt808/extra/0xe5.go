package extra

import (
    "github.com/rayjay214/parser/jt808/errors"
)

type Extra_0xe5 struct {
    serialized []byte
    value      byte
}

func (Extra_0xe5) ID() byte {
    return byte(TypeExtra_0xe5)
}

func (extra Extra_0xe5) Data() []byte {
    return extra.serialized
}

func (extra Extra_0xe5) Value() interface{} {
    return extra.value
}

func (extra Extra_0xe5) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "运动状态"

    if extra.value == 0 {
        m["value"] = "设备静止"
    } else {
        m["value"] = "设备运动"
    }

    return m
}

func (extra *Extra_0xe5) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidExtraLength
    }
    extra.value = data[0]
    extra.serialized = data
    return len(data), nil
}
