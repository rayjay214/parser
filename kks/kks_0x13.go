package kks

import (
	"encoding/json"
	"fmt"
	"parser/common"
	_ "strconv"
)

type Kks_0x13 struct {
    Proto   uint8
    DevInfo uint8
    Voltage uint8
    GSM     uint8
    Lang    uint16
}

func (entity *Kks_0x13) MsgID() MsgID {
    return Msg_0x13
}

func (entity *Kks_0x13) Encode() ([]byte, error) {
    writer := common.NewWriter()

    writer.WriteByte(entity.Proto)

    writer.WriteByte(entity.DevInfo)

    writer.WriteByte(entity.Voltage)

    writer.WriteByte(entity.GSM)

    writer.WriteUint16(entity.Lang)

    return writer.Bytes(), nil
}

func (entity *Kks_0x13) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.Proto, err = reader.ReadByte()
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

    entity.Lang, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}

func (entity Kks_0x13) MarshalJSON() ([]byte, error) {
    type Alias Kks_0x13

    type New0x01 struct {
        Proto string
        Alias
        DevInfo map[string]interface{}
        Voltage string
        GSM     string
        Lang    string
    }

    s := New0x01{
        Alias:   Alias(entity),
        DevInfo: map[string]interface{}{},
    }

    s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

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

    if GetBit(int(entity.Lang), 0) == 1 {
        s.Lang = "中文"
    } else {
        s.Lang = "英文"
    }

    return json.Marshal(s)
}
