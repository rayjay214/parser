package extra

import (
    "github.com/rayjay214/parser/jt808/errors"
    "encoding/hex"
)

type Extra_0xeb struct {
    serialized []byte
    //todo
}

func (Extra_0xeb) ID() byte {
    return byte(TypeExtra_0xeb)
}

func (extra Extra_0xeb) Data() []byte {
    return extra.serialized
}

func (extra Extra_0xeb) Value() interface{} {
    return nil
}

func (extra Extra_0xeb) ToPrint() interface{} {
    m := map[string]interface{}{}
    m["desc"] = "基站信息或扩展协议"
    m["value"] = hex.EncodeToString(extra.serialized)

    return m
}

func (extra *Extra_0xeb) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidExtraLength
    }

    extra.serialized = data
    return len(data), nil
}
