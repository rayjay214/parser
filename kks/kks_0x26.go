package kks

import (
    "encoding/json"
    "fmt"
    "github.com/shopspring/decimal"
    "parser/common"
    _ "strconv"
    "time"
)

type Kks_0x26 struct {
    Proto           uint8
    Time            time.Time
    Satellite       uint8 `json:"-"`
    Lat             decimal.Decimal
    Lng             decimal.Decimal
    Speed           uint8
    StatusDirection uint16
    Mcc             uint16
    Mnc             uint8
    Lac             uint16
    CellId          uint16
    DevInfo         uint8
    Voltage         uint8
    GSM             uint8
    Alarm           uint8
    Lang            uint8
}

func (entity *Kks_0x26) MsgID() MsgID {
    return Msg_0x26
}

func (entity *Kks_0x26) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *Kks_0x26) Decode(data []byte) (int, error) {
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

    entity.StatusDirection, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    status := (entity.StatusDirection & 0xFC00) >> 10
    //route := statusRoute & 0x03FF

    var south, west bool
    if GetBit(int(status), 0) == 0 {
        south = true
    }
    if GetBit(int(status), 1) == 1 {
        west = true
    }

    fmt.Println(status, south, west)

    entity.Lat, entity.Lng = getGeoPoint(latitude, south, longitude, west)

    reader.ReadByte() //len useless

    entity.Mcc, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    entity.Mnc, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Lac, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    reader.ReadByte() //useless

    entity.CellId, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    entity.DevInfo, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Voltage, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.GSM, err = reader.ReadByte()
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

func (entity Kks_0x26) MarshalJSON() ([]byte, error) {
    type Alias Kks_0x26

    type New0x26 struct {
        Proto string
        Alias
        Status    map[string]interface{}
        Direction uint16
        DevInfo   map[string]interface{}
        Voltage   string
        GSM       string
        Alarm     string
        Lang      string
    }

    s := New0x26{
        Alias:   Alias(entity),
        Status:  map[string]interface{}{},
        DevInfo: map[string]interface{}{},
    }

    s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

    status := (entity.StatusDirection & 0xFC00) >> 10
    route := entity.StatusDirection & 0x03FF

    statusMap := map[string]interface{}{}

    if GetBit(int(status), 2) == 1 {
        statusMap["定位状态"] = "GPS已定位"
    } else {
        statusMap["定位状态"] = "GPS未定位"
    }
    if GetBit(int(status), 3) == 1 {
        statusMap["定位类型"] = "差分定位"
    } else {
        statusMap["定位类型"] = "GPS实时"
    }
    if GetBit(int(status), 4) == 1 {
        statusMap["设备状态"] = "运动"
    } else {
        statusMap["设备状态"] = "静止"
    }
    s.Status = statusMap
    s.Direction = route

    devInfoMap := map[string]interface{}{}
    if GetBit(int(entity.DevInfo), 0) == 1 {
        devInfoMap["设防状态"] = "设防"
    } else {
        devInfoMap["设防状态"] = "撤防"
    }
    if GetBit(int(entity.DevInfo), 1) == 1 {
        devInfoMap["ACC"] = "高"
    } else {
        devInfoMap["ACC"] = "低"
    }
    if GetBit(int(entity.DevInfo), 2) == 1 {
        devInfoMap["电源状态"] = "已接电源充电"
    } else {
        devInfoMap["电源状态"] = "未接电源充电"
    }
    if GetBit(int(entity.DevInfo), 6) == 1 {
        devInfoMap["定位状态"] = "GPS已定位"
    } else {
        devInfoMap["定位状态"] = "GPS未定位"
    }
    if GetBit(int(entity.DevInfo), 7) == 1 {
        devInfoMap["油电状态"] = "油电断开"
    } else {
        devInfoMap["油电状态"] = "油电接通"
    }
    s.DevInfo = devInfoMap

    switch entity.Voltage {
    case 0x00:
        s.Voltage = "无电"
    case 0x01:
        s.Voltage = "电量极低"
    case 0x02:
        s.Voltage = "电量很低"
    case 0x03:
        s.Voltage = "电量低"
    case 0x04:
        s.Voltage = "电量中"
    case 0x05:
        s.Voltage = "电量高"
    case 0x06:
        s.Voltage = "电量极高"
    }

    switch entity.GSM {
    case 0x00:
        s.GSM = "无信号"
    case 0x01:
        s.GSM = "信号极弱"
    case 0x02:
        s.GSM = "信号较弱"
    case 0x03:
        s.GSM = "信号良好"
    case 0x04:
        s.GSM = "信号强"
    }

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

    if entity.Lang == 0x01 {
        s.Lang = "中文回复"
    } else if entity.Lang == 0x02 {
        s.Lang = "英文回复"
    } else {
        s.Lang = "不需要回复"
    }

    return json.Marshal(s)
}
