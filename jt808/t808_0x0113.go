package jt808

import (
    "encoding/json"
    "github.com/rayjay214/parser/common"
    "github.com/rayjay214/parser/jt808/errors"
)

// 上传蓝牙相关定位模式
type T808_0x0113 struct {
    Mode               byte
    IsReportInInterval byte `json:"是否实时上报"`
    IsKeepGprs         byte `json:"是否保持GSM功能"`
    IsReportInTime     byte `json:"是否等待时间"`
    IsWifiInPriority   byte `json:"是否WiFi优先"`
    GprsInterval       uint32
    ConnInterval       uint32
}

func (entity T808_0x0113) MarshalJSON() ([]byte, error) {
    type Alias T808_0x0113

    type NewT808_0x0113 struct {
        Alias
        Mode string
    }

    s := NewT808_0x0113{
        Alias: Alias(entity),
    }

    switch entity.Mode {
    case 15:
        s.Mode = "鹅卵石自由模式"
    case 16:
        s.Mode = "鹅卵石待机模式"
    case 17:
        s.Mode = "鹅卵石定时模式"
    case 18:
        s.Mode = "鹅卵石实时模式"
    default:
        s.Mode = "其他模式"
    }

    return json.Marshal(s)
}

func (entity *T808_0x0113) MsgID() MsgID {
    return MsgT808_0x0113
}

func (entity *T808_0x0113) Encode() ([]byte, error) {
    return nil, nil
}

func (entity *T808_0x0113) Decode(data []byte) (int, error) {
    if len(data) < 6 {
        return 0, errors.ErrInvalidBody
    }
    reader := common.NewReader(data)

    var err error

    entity.Mode, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.IsReportInInterval, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.IsKeepGprs, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.IsReportInTime, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.IsWifiInPriority, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.GprsInterval, err = reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    entity.ConnInterval, err = reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    return 0, nil
}
