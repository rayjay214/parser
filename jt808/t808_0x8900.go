package jt808

import (
    "github.com/rayjay214/parser/common"
    "github.com/rayjay214/parser/jt808/errors"
)

// 数据下行透传
type T808_0x8900 struct {
    // 透传消息类型
    Type byte
    // 透传消息内容
    Data string
}

func (entity *T808_0x8900) MsgID() MsgID {
    return MsgT808_0x8900
}

func (entity *T808_0x8900) Encode() ([]byte, error) {
    writer := common.NewWriter()
    writer.WriteByte(entity.Type)
    writer.WriteString(entity.Data)
    return writer.Bytes(), nil
}

func (entity *T808_0x8900) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidBody
    }
    entity.Type = data[0]
    reader := common.NewReader(data[1:])
    if reader.Len() > 0 {
        data, err := reader.ReadString()
        if err != nil {
            return 0, err
        }
        entity.Data = data
    }
    return len(data), nil
}
