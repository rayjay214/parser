package jt808

import (
    "github.com/rayjay214/parser/common"
)

type T808_0x0106 struct {
    Result byte
}

func (entity *T808_0x0106) MsgID() MsgID {
    return MsgT808_0x0106
}

func (entity *T808_0x0106) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *T808_0x0106) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    entity.Result, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
