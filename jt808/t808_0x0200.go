package jt808

import (
    _ "encoding/hex"
    "encoding/json"
    "fmt"
    "github.com/shopspring/decimal"
    log "github.com/sirupsen/logrus"
    "math"
    "github.com/rayjay214/parser/common"
    "github.com/rayjay214/parser/jt808/errors"
    "github.com/rayjay214/parser/jt808/extra"
    "reflect"
    "time"
)

// 纬度类型
type LatitudeType int

const (
    _ LatitudeType = iota
    // 北纬
    NorthLatitudeType = 0
    // 南纬
    SouthLatitudeType = 1
)

// 经度类型
type LongitudeType int

const (
    _ LongitudeType = iota
    // 东经
    EastLongitudeType = 0
    // 西经
    WestLongitudeType = 1
)

// 位置状态
type T808_0x0200_Status uint32

// 获取Acc状态
func (status T808_0x0200_Status) GetAccState() bool {
    return GetBitUint32(uint32(status), 0)
}

// 是否正在定位
func (status T808_0x0200_Status) Positioning() bool {
    return GetBitUint32(uint32(status), 1)
}

// 设置南纬
func (status *T808_0x0200_Status) SetSouthLatitude(b bool) {
    SetBitUint32((*uint32)(status), 2, b)
}

// 设置西经
func (status *T808_0x0200_Status) SetWestLongitude(b bool) {
    SetBitUint32((*uint32)(status), 3, b)
}

// 获取纬度类型
func (status T808_0x0200_Status) GetLatitudeType() LatitudeType {
    if GetBitUint32(uint32(status), 2) {
        return SouthLatitudeType
    }
    return NorthLatitudeType
}

// 获取经度类型
func (status T808_0x0200_Status) GetLongitudeType() LongitudeType {
    if GetBitUint32(uint32(status), 3) {
        return WestLongitudeType
    }
    return EastLongitudeType
}

// 汇报位置
type T808_0x0200 struct {
    // 警告
    Alarm uint32
    // 状态
    Status T808_0x0200_Status
    // 纬度
    Lat decimal.Decimal
    // 经度
    Lng decimal.Decimal
    // 海拔高度
    // 单位：米
    Altitude uint16
    // 速度
    // 单位：1/10km/h
    Speed uint16
    // 方向
    // 0-359，正北为 0，顺时针
    Direction uint16
    // 时间
    Time time.Time
    // 附加信息
    Extras []extra.Entity
}

func (entity *T808_0x0200) MsgID() MsgID {
    return MsgT808_0x0200
}

func (entity *T808_0x0200) Encode() ([]byte, error) {
    writer := common.NewWriter()

    // 写入警告标志
    writer.WriteUint32(entity.Alarm)

    // 计算经纬度
    mul := decimal.NewFromFloat(1000000)
    lat := entity.Lat.Mul(mul).IntPart()
    if lat < 0 {
        entity.Status.SetSouthLatitude(true)
    }
    lng := entity.Lng.Mul(mul).IntPart()
    if lng < 0 {
        entity.Status.SetWestLongitude(true)
    }

    // 写入状态信息
    writer.WriteUint32(uint32(entity.Status))

    // 写入纬度信息
    writer.WriteUint32(uint32(math.Abs(float64(lat))))

    // 写入经度信息
    writer.WriteUint32(uint32(math.Abs(float64(lng))))

    // 写入海拔高度
    writer.WriteUint16(entity.Altitude)

    // 写入速度信息
    writer.WriteUint16(entity.Speed)

    // 写入方向信息
    writer.WriteUint16(entity.Direction)

    // 写入时间信息
    writer.WriteBcdTime(entity.Time)

    // 写入附加信息
    for i := 0; i < len(entity.Extras); i++ {
        ext := entity.Extras[i]
        if ext == nil || reflect.ValueOf(ext).IsNil() {
            continue
        }
        data := ext.Data()
        full := make([]byte, len(data)+2)
        full[0], full[1] = ext.ID(), byte(len(data))
        copy(full[2:], data)
        writer.Write(full)
    }
    return writer.Bytes(), nil
}

func (entity *T808_0x0200) Decode(data []byte) (int, error) {
    if len(data) < 28 {
        return 0, errors.ErrInvalidBody
    }
    reader := common.NewReader(data)

    // 读取警告标志
    var err error
    entity.Alarm, err = reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    // 读取状态信息
    status, err := reader.ReadUint32()
    if err != nil {
        return 0, err
    }
    entity.Status = T808_0x0200_Status(status)

    // 读取纬度信息
    latitude, err := reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    // 读取经度信息
    longitude, err := reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    entity.Lat, entity.Lng = common.GetGeoPoint(
        latitude, entity.Status.GetLatitudeType() == SouthLatitudeType,
        longitude, entity.Status.GetLongitudeType() == WestLongitudeType)

    // 读取海拔高度
    entity.Altitude, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    // 读取行驶速度
    entity.Speed, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    // 读取行驶方向
    entity.Direction, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    // 读取上报时间
    entity.Time, err = reader.ReadBcdTime()
    if err != nil {
        return 0, err
    }

    // 解码附加信息
    extras := make([]extra.Entity, 0)
    buffer := data[len(data)-reader.Len():]
    for {
        if len(buffer) < 2 {
            break
        }
        id, length := buffer[0], int(buffer[1])
        buffer = buffer[2:]
        if len(buffer) < length {
            return 0, errors.ErrInvalidExtraLength
        }

        extraEntity, count, err := extra.Decode(id, buffer[:length])
        if err != nil {
            if err == errors.ErrTypeNotRegistered {
                buffer = buffer[length:]
                log.WithFields(log.Fields{
                    "id": fmt.Sprintf("0x%x", id),
                }).Warn("[JT/T808] unknown T808_0x0200 extra type")
                continue
            }
            return 0, err
        }
        if count != length {
            return 0, errors.ErrInvalidExtraLength
        }
        extras = append(extras, extraEntity)
        buffer = buffer[length:]
    }
    if len(extras) > 0 {
        entity.Extras = extras
    }
    return len(data) - reader.Len(), nil
}

func (entity T808_0x0200) MarshalJSON() ([]byte, error) {
    type Alias T808_0x0200

    type New0200 struct {
        Alias
        Alarm   []string
        Status  map[string]interface{}
        Extras  map[string]interface{}
        LocType string
        Speed   float32
    }

    s := New0200{
        Alias:  Alias(entity),
        Status: map[string]interface{}{},
        Extras: map[string]interface{}{},
    }

    s.Speed = float32(entity.Speed) / 10

    var isLocated bool
    var isMoving bool
    var isWifi bool

    fGetAlarmbit := func(value uint32, offset int) uint32 {
        return uint32(value) & (1 << offset) >> offset
    }

    s.Alarm = make([]string, 0)
    if fGetAlarmbit(entity.Alarm, 0) == 1 {
        s.Alarm = append(s.Alarm, "sos报警")
    }
    if fGetAlarmbit(entity.Alarm, 1) == 1 {
        s.Alarm = append(s.Alarm, "超速报警")
    }
    if fGetAlarmbit(entity.Alarm, 2) == 1 {
        s.Alarm = append(s.Alarm, "解除超速报警")
    }
    if fGetAlarmbit(entity.Alarm, 5) == 1 {
        s.Alarm = append(s.Alarm, "GPS信号弱")
    }
    if fGetAlarmbit(entity.Alarm, 6) == 1 {
        s.Alarm = append(s.Alarm, "正常开机报警")
    }
    if fGetAlarmbit(entity.Alarm, 7) == 1 {
        s.Alarm = append(s.Alarm, "终端主电源欠压")
    }
    if fGetAlarmbit(entity.Alarm, 8) == 1 {
        s.Alarm = append(s.Alarm, "防拆报警")
    }
    if fGetAlarmbit(entity.Alarm, 9) == 1 {
        s.Alarm = append(s.Alarm, "震动报警")
    }
    if fGetAlarmbit(entity.Alarm, 10) == 1 {
        s.Alarm = append(s.Alarm, "开机报警")
    }
    if fGetAlarmbit(entity.Alarm, 11) == 1 {
        s.Alarm = append(s.Alarm, "关机报警")
    }
    if fGetAlarmbit(entity.Alarm, 12) == 1 {
        s.Alarm = append(s.Alarm, "光感脱落报警")
    }
    if fGetAlarmbit(entity.Alarm, 13) == 1 {
        s.Alarm = append(s.Alarm, "光感恢复报警")
    }
    if fGetAlarmbit(entity.Alarm, 14) == 1 {
        s.Alarm = append(s.Alarm, "非法位移报警")
    }
    if fGetAlarmbit(entity.Alarm, 15) == 1 {
        s.Alarm = append(s.Alarm, "强烈震动报警")
    }
    if fGetAlarmbit(entity.Alarm, 16) == 1 {
        s.Alarm = append(s.Alarm, "解除防拆报警或怠速告警")
    }
    if fGetAlarmbit(entity.Alarm, 17) == 1 {
        s.Alarm = append(s.Alarm, "拆除告警")
    }

    fGetbit := func(value T808_0x0200_Status, offset int) uint32 {
        return uint32(value) & (1 << offset) >> offset
    }

    statusMap := map[string]interface{}{}
    if fGetbit(entity.Status, 0) == 1 {
        statusMap["acc"] = "开启"
    } else {
        statusMap["acc"] = "关闭"
    }
    if fGetbit(entity.Status, 1) == 1 {
        statusMap["定位状态"] = "已定位"
        isLocated = true
    } else {
        statusMap["定位状态"] = "未定位"
        isLocated = false
    }
    if fGetbit(entity.Status, 2) == 1 {
        statusMap["纬度"] = "南纬"
    } else {
        statusMap["纬度"] = "北纬"
    }
    if fGetbit(entity.Status, 3) == 1 {
        statusMap["经度"] = "西经"
    } else {
        statusMap["经度"] = "东经"
    }
    if fGetbit(entity.Status, 6) == 1 {
        statusMap["设防状态"] = "设防"
    } else {
        statusMap["设防状态"] = "撤防"
    }
    if fGetbit(entity.Status, 10) == 1 {
        statusMap["油路状态"] = "断开"
    } else {
        statusMap["油路状态"] = "正常"
    }
    if fGetbit(entity.Status, 11) == 1 {
        statusMap["主电源"] = "断开"
    } else {
        statusMap["主电源"] = "接通"
    }
    if fGetbit(entity.Status, 18) == 1 {
        statusMap["GPS"] = "使用"
    } else {
        statusMap["GPS"] = "未使用"
    }
    if fGetbit(entity.Status, 19) == 1 {
        statusMap["北斗"] = "使用"
    } else {
        statusMap["北斗"] = "未使用"
    }
    if fGetbit(entity.Status, 20) == 1 {
        statusMap["GLONASS"] = "使用"
    } else {
        statusMap["GLONASS"] = "未使用"
    }
    if fGetbit(entity.Status, 21) == 1 {
        statusMap["Galileo"] = "使用"
    } else {
        statusMap["Galileo"] = "未使用"
    }
    if fGetbit(entity.Status, 22) == 1 {
        statusMap["车辆状态"] = "行驶"
        isMoving = true
    } else {
        statusMap["车辆状态"] = "静止"
        isMoving = false
    }
    s.Status = statusMap

    isLbs := false
    for _, v := range entity.Extras {
        if v.ID() == 0xec || v.ID() == 0x54 {
            isWifi = true
        }

        if v.ID() == 0x5d || v.ID() == 0xe1 || v.ID() == 0xee || v.ID() == 0xeb {
            isLbs = true
        }

        if v.ID() == 0x05 {
            if v.Value() == byte(1) {
                isMoving = true
            } else {
                isMoving = false
            }
        }

        if v.ID() == 0xe5 {
            if v.Value() == byte(1) {
                isMoving = true
            } else {
                isMoving = false
            }
        }

        if v.ID() == 0xe7 {
            e7, ok := v.Value().(extra.Extra_0xe7_Value)
            if ok {
                if e7.ShakeAlarm == byte(1) {
                    s.Alarm = append(s.Alarm, "震动告警")
                }
            }
        }

        var strId string
        var val interface{}
        strId = "0x" + fmt.Sprintf("%02x", v.ID())
        switch vPrint := v.ToPrint().(type) {
        case map[string]interface{}:
            val = vPrint
        default:
            fmt.Println("%T", vPrint)
        }
        s.Extras[strId] = val
    }

    if isLocated {
        if isMoving {
            s.LocType = "运动GPS"
        } else {
            s.LocType = "静止GPS"
        }
    } else {
        if isMoving {
            if isWifi && isLbs {
                s.LocType = "运动混合定位"
            } else if isWifi {
                s.LocType = "运动WIFI"
            } else {
                s.LocType = "运动基站"
            }

        } else {
            if isWifi && isLbs {
                s.LocType = "静止混合定位"
            } else if isWifi {
                s.LocType = "静止WIFI"
            } else {
                s.LocType = "静止基站"
            }
        }
    }
    return json.Marshal(s)
}
