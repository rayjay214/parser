package jt808

import (
    "parser/common"
)

// 请求同步时间
type T808_0x1107 struct {
    IccidLen      byte
    Iccid         string
    DeviceTypeLen byte
    DeviceType    string
    VersionLen    byte
    Version       string
}

func (entity *T808_0x1107) MsgID() MsgID {
    return MsgT808_0x1107
}

func (entity *T808_0x1107) Encode() ([]byte, error) {
    return nil, nil
}

func (entity *T808_0x1107) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    entity.IccidLen, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Iccid, err = reader.ReadString(20)
    if err != nil {
        return 0, err
    }

    entity.DeviceTypeLen, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.DeviceType, err = reader.ReadString(int(entity.DeviceTypeLen))
    if err != nil {
        return 0, err
    }

    entity.VersionLen, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Version, err = reader.ReadString(int(entity.VersionLen))
    if err != nil {
        return 0, err
    }

    return 0, nil
}
