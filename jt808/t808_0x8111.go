package jt808

import (
    "parser/common"
    "fmt"
)

// 终端应答
type T808_0x8111 struct {
    LocType  byte
    LocTime  string
    PlateLen byte
    Plate    string
    AddrLen  uint16
    Addr     string
}

func (entity *T808_0x8111) MsgID() MsgID {
    return MsgT808_0x8111
}

func (entity *T808_0x8111) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *T808_0x8111) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    entity.LocType, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    year, _ := reader.ReadByte()
    month, _ := reader.ReadByte()
    day, _ := reader.ReadByte()
    hour, _ := reader.ReadByte()
    minute, _ := reader.ReadByte()
    second, _ := reader.ReadByte()
    entity.LocTime = fmt.Sprintf("20%02d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)

    entity.PlateLen, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    Plate, err := reader.Read(int(entity.PlateLen))
    if err != nil {
        return 0, err
    }
    entity.Plate = common.BytesToString(Plate[:])

    entity.AddrLen, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    Addr, err := reader.Read(int(entity.AddrLen))
    if err != nil {
        return 0, err
    }
    entity.Addr = common.BytesToString(Addr[:])

    return len(data) - reader.Len(), nil
}
