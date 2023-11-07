package jt808

import (
    "parser/common"
    "parser/jt808/errors"
)

// 终端鉴权
type T808_0x0102 struct {
    // 鉴权码
    AuthKey string
}

func (entity *T808_0x0102) MsgID() MsgID {
    return MsgT808_0x0102
}

func (entity *T808_0x0102) Encode() ([]byte, error) {
    writer := common.NewWriter()
    if len(entity.AuthKey) > 0 {
        if err := writer.WriteString(entity.AuthKey); err != nil {
            return nil, err
        }
    }
    return writer.Bytes(), nil
}

func (entity *T808_0x0102) Decode(data []byte) (int, error) {
    if len(data) == 0 {
        return 0, errors.ErrInvalidBody
    }

    reader := common.NewReader(data)
    authKey, err := reader.ReadString()
    if err != nil {
        return 0, err
    }
    entity.AuthKey = authKey
    return len(data) - reader.Len(), nil
}
