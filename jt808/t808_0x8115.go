package jt808

import (
    "github.com/rayjay214/parser/common"
)

// 终端应答
type T808_0x8115 struct {
    SessionId string
}

func (entity *T808_0x8115) MsgID() MsgID {
    return MsgT808_0x8115
}

func (entity *T808_0x8115) Encode() ([]byte, error) {
    writer := common.NewWriter()

    writer.WriteString(entity.SessionId)

    return writer.Bytes(), nil
}

func (entity *T808_0x8115) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    entity.SessionId, err = reader.ReadString(8)
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
