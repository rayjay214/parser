package extra

import (
    "github.com/rayjay214/parser/jt808/errors"
    "encoding/binary"
    "encoding/json"
)

// 基站信息
type Extra_0xe7 struct {
    serialized []byte
    value      Extra_0xe7_Value
}

type Extra_0xe7_Value struct {
    ShakeAlarm    byte `json:"-"`
    FenceStatus   byte `json:"-"`
    SleepStatus   byte `json:"-"`
    SleepCheckWay byte `json:"-"`
}

func (entity Extra_0xe7_Value) MarshalJSON() ([]byte, error) {
    type Alias Extra_0xe7_Value

    type NewExtra_0xe7_Value struct {
        Alias
        ShakeAlarm    string `json:"震动报警"`
        FenceStatus   string `json:"设防状态"`
        SleepStatus   string `json:"休眠状态"`
        SleepCheckWay string `json:"休眠判断途径"`
    }

    s := NewExtra_0xe7_Value{
        Alias: Alias(entity),
    }

    if entity.ShakeAlarm == 1 {
        s.ShakeAlarm = "报警"
    } else {
        s.ShakeAlarm = "未报警"
    }
    if entity.FenceStatus == 1 {
        s.FenceStatus = "设防"
    } else {
        s.FenceStatus = "未设防"
    }
    if entity.SleepStatus == 1 {
        s.SleepStatus = "休眠"
    } else {
        s.SleepStatus = "唤醒"
    }
    if entity.SleepCheckWay == 1 {
        s.SleepCheckWay = "根据E7"
    } else {
        s.SleepCheckWay = "根据ACC"
    }

    return json.Marshal(s)
}

func (Extra_0xe7) ID() byte {
    return byte(TypeExtra_0xe7)
}

func (extra Extra_0xe7) Data() []byte {
    return extra.serialized
}

func (extra Extra_0xe7) Value() interface{} {
    return extra.value
}

func (extra Extra_0xe7) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "状态扩展"
    m["value"] = extra.Value()

    return m
}

func (extra *Extra_0xe7) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidExtraLength
    }
    alarm := binary.BigEndian.Uint16(data)
    status := binary.BigEndian.Uint16(data[2:])

    extra.value.ShakeAlarm = byte(alarm & 1)
    extra.value.FenceStatus = byte(status & 1)
    extra.value.SleepStatus = byte(status & (1 << 1) >> 1)
    extra.value.SleepCheckWay = byte(status & (1 << 2) >> 2)
    extra.serialized = data
    return len(data), nil
}
