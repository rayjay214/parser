package jt808

import (
    "parser/common"
)

// 终端应答
type T808_0x0117 struct {
    PkgSize byte
    PkgNo   byte
}

func (entity *T808_0x0117) MsgID() MsgID {
    return MsgT808_0x0117
}

func (entity *T808_0x0117) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *T808_0x0117) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    entity.PkgSize, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.PkgNo, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
