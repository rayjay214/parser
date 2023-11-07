package extra

import (
    _ "parser/jt808/errors"
    "encoding/hex"
    "parser/common"
    "encoding/json"
)

type Extra_0xf7 struct {
    serialized []byte
    value      Extra_0xf7_Value
}

type Extra_0xf7_Value struct {
    Mcc           []byte
    Mnc           []byte
    LenApn        byte `json:"-"`
    Apn           string
    LenApnAccount byte `json:"-"`
    ApnAccount    string
    LenApnPwd     byte `json:"-"`
    ApnPwd        string
}

func (entity Extra_0xf7_Value) MarshalJSON() ([]byte, error) {
    type Alias Extra_0xf7_Value

    type NewExtra_0xf7_Value struct {
        Alias
        Mcc string
        Mnc string
    }

    s := NewExtra_0xf7_Value{
        Alias: Alias(entity),
    }

    s.Mcc = hex.EncodeToString(entity.Mcc)
    s.Mnc = hex.EncodeToString(entity.Mnc)

    return json.Marshal(s)
}

func (Extra_0xf7) ID() byte {
    return byte(TypeExtra_0xf7)
}

func (extra Extra_0xf7) Data() []byte {
    return extra.serialized
}

func (extra Extra_0xf7) Value() interface{} {
    return extra.value
}

func (extra Extra_0xf7) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "IMEI上传"
    m["value"] = extra.Value()

    return m
}

func (extra *Extra_0xf7) Decode(data []byte) (int, error) {

    extra.serialized = data

    reader := common.NewReader(data)

    extra.value.Mcc, _ = reader.Read(3)
    extra.value.Mnc, _ = reader.Read(3)
    extra.value.LenApn, _ = reader.ReadByte()
    extra.value.Apn, _ = reader.ReadString(int(extra.value.LenApn))
    extra.value.LenApnAccount, _ = reader.ReadByte()
    if extra.value.LenApnAccount > 0 {
        extra.value.ApnAccount, _ = reader.ReadString(int(extra.value.LenApnAccount))
    }
    extra.value.LenApnPwd, _ = reader.ReadByte()
    if extra.value.LenApnPwd > 0 {
        extra.value.ApnPwd, _ = reader.ReadString(int(extra.value.LenApnPwd))
    }

    return len(data), nil
}
