package kks

import (
    "encoding/json"
    "fmt"
    "github.com/shopspring/decimal"
    "github.com/rayjay214/parser/common"
    _ "strconv"
    "time"
)

type Kks_0x2a struct {
    Proto           uint8
    Time            time.Time
    Satellite       uint8 `json:"-"`
    Lat             decimal.Decimal
    Lng             decimal.Decimal
    Speed           uint8
    StatusDirection uint16
    Phone           string
    Alarm           uint8
    Lang            uint8
}

func (entity *Kks_0x2a) MsgID() MsgID {
    return Msg_0x2a
}

func (entity *Kks_0x2a) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *Kks_0x2a) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.Proto, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Time, err = reader.ReadStrTime()
    if err != nil {
        return 0, err
    }

    entity.Satellite, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    latitude, err := reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    longitude, err := reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    entity.Speed, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    statusRoute, err := reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    status := (statusRoute & 0xFC00) >> 10
    //route := statusRoute & 0x03FF

    entity.StatusDirection = statusRoute

    var south, west bool
    if GetBit(int(status), 0) == 0 {
        south = true
    }
    if GetBit(int(status), 1) == 1 {
        west = true
    }

    entity.Lat, entity.Lng = getGeoPoint(latitude, south, longitude, west)

    entity.Phone, err = reader.ReadString(21)
    if err != nil {
        return 0, err
    }

    entity.Alarm, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Lang, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}

func (entity Kks_0x2a) MarshalJSON() ([]byte, error) {
    type Alias Kks_0x2a

    type New0x2a struct {
        Proto string
        Alarm string
        Lang  string
        Alias
    }

    s := New0x2a{
        Alias: Alias(entity),
    }

    s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

    fmt.Println(entity.Alarm)
    switch entity.Alarm {
    case 0x00:
        s.Alarm = "正常"
    case 0x01:
        s.Alarm = "SOS求救"
    case 0x02:
        s.Alarm = "断电报警"
    case 0x03:
        s.Alarm = "震动报警"
    case 0x04:
        s.Alarm = "进围栏报警"
    case 0x05:
        s.Alarm = "出围栏报警"
    case 0x06:
        s.Alarm = "超速报警"
    case 0x09:
        s.Alarm = "位移报警"
    case 0x0a:
        s.Alarm = "进GPS盲区报警"
    case 0x0b:
        s.Alarm = "出GPS盲区报警"
    case 0x0c:
        s.Alarm = "开机报警"
    case 0x0d:
        s.Alarm = "GPS 第一次定位报警"
    case 0x0e:
        s.Alarm = "外电低电报警"
    case 0x0f:
        s.Alarm = "外电低电保护报警"
    case 0x10:
        s.Alarm = "换卡报警"
    case 0x11:
        s.Alarm = "低电关机报警"
    case 0x12:
        s.Alarm = "外电低电保护后飞行模式报警"
    case 0x13:
        s.Alarm = "拆卸报警"
    case 0x14:
        s.Alarm = "门报警"
    case 0x15:
        s.Alarm = "低电关机报警"
    case 0x16:
        s.Alarm = "声控报警"
    case 0x17:
        s.Alarm = "伪基站报警"
    case 0xFF:
        s.Alarm = "ACC关"
    case 0xFE:
        s.Alarm = "ACC开"
    }

    if GetBit(int(entity.Lang), 0) == 1 {
        s.Lang = "中文"
    } else {
        s.Lang = "英文"
    }

    return json.Marshal(s)
}
