package jt808

import (
    "github.com/rayjay214/parser/common"
)

// 终端应答
type T808_0x0116 struct {
    RecordStatus byte
}

func (entity *T808_0x0116) MsgID() MsgID {
    return MsgT808_0x0116
}

func (entity *T808_0x0116) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *T808_0x0116) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    entity.RecordStatus, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
