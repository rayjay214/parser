package jt808

import (
    "github.com/rayjay214/parser/common"
    "fmt"
)

type T808_0x8131 struct {
    Count   byte     //时间段组数
    Periods []string //时间段
    Mode    byte     //定时开关机模式
    DayCnt  byte     //天数
    Days    []string //日期
}

func (entity *T808_0x8131) MsgID() MsgID {
    return MsgT808_0x8110
}

func (entity *T808_0x8131) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *T808_0x8131) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    entity.Count, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Periods = make([]string, 0, entity.Count)
    for i := 0; i < int(entity.Count); i++ {
        temp, err := reader.Read(4)
        if err != nil {
            return 0, err
        }
        period := fmt.Sprintf("%02x:%02x -- %02x:%02x", temp[0], temp[1], temp[2], temp[3])

        entity.Periods = append(entity.Periods, period)
    }

    entity.Mode, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    if entity.Mode != 0 {
        entity.DayCnt, err = reader.ReadByte()
        if err != nil {
            return 0, err
        }

        entity.Days = make([]string, 0, entity.DayCnt)
        for i := 0; i < int(entity.DayCnt); i++ {
            temp, err := reader.Read(1)
            if err != nil {
                return 0, err
            }
            day := fmt.Sprintf("%02x", temp[0])
            entity.Days = append(entity.Days, day)
        }

    }

    return len(data) - reader.Len(), nil
}
