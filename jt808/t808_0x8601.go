package jt808

import (
    "github.com/rayjay214/parser/common"
    "github.com/rayjay214/parser/jt808/errors"
)

// 删除圆形区域
type T808_0x8601 struct {
    // 区域ID列表
    IDs []uint32
}

func (entity *T808_0x8601) MsgID() MsgID {
    return MsgT808_0x8601
}

func (entity *T808_0x8601) Encode() ([]byte, error) {
    writer := common.NewWriter()

    // 写入ID总数
    writer.WriteByte(byte(len(entity.IDs)))

    // 写入ID列表
    for _, id := range entity.IDs {
        writer.WriteUint32(id)
    }
    return writer.Bytes(), nil
}

func (entity *T808_0x8601) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidBody
    }

    count := int(data[0])
    reader := common.NewReader(data[1:])
    entity.IDs = make([]uint32, 0, count)
    for i := 0; i < count; i++ {
        id, err := reader.ReadUint32()
        if err != nil {
            return 0, err
        }
        entity.IDs = append(entity.IDs, id)
    }
    return len(data) - reader.Len() - 1, nil
}
