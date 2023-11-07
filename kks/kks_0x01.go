package kks

import (
    "encoding/json"
    "fmt"
    "github.com/rayjay214/parser/common"
    "strconv"
)

type Kks_0x01 struct {
    Proto        uint8
    Imei         uint64
    Code         uint16 `json:"-"`
    TimezoneLang uint16 `json:"-"`
}

func (entity *Kks_0x01) MsgID() MsgID {
    return Msg_0x01
}

func (entity *Kks_0x01) Encode() ([]byte, error) {
    writer := common.NewWriter()

    writer.WriteByte(entity.Proto)

    writer.Write(stringToBCD(strconv.FormatUint(entity.Imei, 10), 8))

    writer.WriteUint16(entity.Code)

    writer.WriteUint16(entity.TimezoneLang)

    return writer.Bytes(), nil
}

func (entity *Kks_0x01) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.Proto, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    // 读取终端号码
    temp, err := reader.Read(8)
    if err != nil {
        return 0, ErrInvalidHeader
    }
    imei, err := strconv.ParseUint(bcdToString(temp), 10, 64)
    if err != nil {
        return 0, err
    }
    entity.Imei = imei

    if len(data) >= 19 {
        entity.Code, err = reader.ReadUint16()
        if err != nil {
            return 0, err
        }

        entity.TimezoneLang, err = reader.ReadUint16()
        if err != nil {
            return 0, err
        }
    } else {
        entity.TimezoneLang = 0xffff
    }

    return len(data) - reader.Len(), nil
}

func (entity Kks_0x01) MarshalJSON() ([]byte, error) {
    type Alias Kks_0x01

    type New0x01 struct {
        Proto string
        Alias
        TimeZone string
    }

    s := New0x01{
        Alias: Alias(entity),
    }

    s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

    timezone := (entity.TimezoneLang & 0xFFF0) >> 4
    status := entity.TimezoneLang & 0x000F

    flag := status & (1 << 3) >> 3
    if flag == 0 {
        s.TimeZone = "东" + strconv.Itoa(int(timezone/100)) + "区"
    } else {
        s.TimeZone = "西" + strconv.Itoa(int(timezone/100)) + "区"
    }

    if entity.TimezoneLang == 0xffff {
        s.TimeZone = "无"
    }

    return json.Marshal(s)
}
