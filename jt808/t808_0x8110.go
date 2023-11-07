package jt808

import (
    "github.com/rayjay214/parser/common"
)

type LoopData struct {
    DayOfWeek   byte     //当前系统星期几
    TimeSize    byte     //时间定位个数
    LocTimeList []string //具体时间
}

// 终端应答
type T808_0x8110 struct {
    SysDayOfWeek byte       //当前系统星期几
    LocType      byte       //定位类型
    LocSize      uint16     //定位包个数
    LoopDataList []LoopData //定位包列表
}

func (entity *T808_0x8110) MsgID() MsgID {
    return MsgT808_0x8110
}

func (entity *T808_0x8110) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *T808_0x8110) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    entity.SysDayOfWeek, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.LocType, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.LocSize, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    entity.LoopDataList = make([]LoopData, 0, entity.LocSize)
    for i := 0; i < int(entity.LocSize); i++ {
        var loopData LoopData
        loopData.DayOfWeek, err = reader.ReadByte()
        if err != nil {
            return 0, err
        }

        loopData.TimeSize, err = reader.ReadByte()
        if err != nil {
            return 0, err
        }

        loopData.LocTimeList = make([]string, 0, loopData.TimeSize)
        for j := 0; j < int(loopData.TimeSize); j++ {
            var locTime string
            locTime, err = reader.ReadString(5)
            if err != nil {
                return 0, err
            }
            loopData.LocTimeList = append(loopData.LocTimeList, locTime)
        }
        entity.LoopDataList = append(entity.LoopDataList, loopData)
    }

    return len(data) - reader.Len(), nil
}
