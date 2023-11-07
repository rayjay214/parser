package ota

import (
	_ "fmt"
	"parser/common"
	"strconv"
)

type Body_0x91 struct {
    UpgradeFlag byte
    UpgradeType byte
    TaskId      uint32
    SrcVersion  string
    DstVersion  string
    UrlLen      byte
    Url         string
    Md5         string
}

type Body_0x11 struct {
    DevType     string
    Sn          uint64
    CurrVersion string
    Lbsinfo     []byte
}

type Body_0x12 struct {
    TaskId      uint32
    UpgradeType byte
    Reason      byte
    CurrVersion string
}

func (entity *Body_0x11) Encode() ([]byte, error) {
    writer := common.NewWriter()

    writer.WriteString(entity.DevType, 10)

    writer.Write(stringToBCD(strconv.FormatUint(entity.Sn, 10), 6))

    writer.WriteString(entity.CurrVersion, 20)

    writer.Write(entity.Lbsinfo, 10)

    return writer.Bytes(), nil
}

func (entity *Body_0x11) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.DevType, err = reader.ReadString(10)
    if err != nil {
        return 0, err
    }

    temp, err := reader.Read(6)
    if err != nil {
        return 0, ErrInvalidHeader
    }
    bcdStr := common.BcdToString(temp)
    if bcdStr == "" {
        entity.Sn = 0
    } else {
        entity.Sn, err = strconv.ParseUint(common.BcdToString(temp), 10, 64)
        if err != nil {
            return 0, err
        }
    }

    entity.CurrVersion, err = reader.ReadString(20)
    if err != nil {
        return 0, err
    }

    entity.Lbsinfo, err = reader.Read(10)
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}

func (entity *Body_0x91) MsgID() MsgID {
    return Msg_0x91
}

func (entity *Body_0x12) Encode() ([]byte, error) {
    writer := common.NewWriter()

    writer.WriteUint32(entity.TaskId)

    writer.WriteByte(entity.UpgradeType)

    writer.WriteByte(entity.Reason)

    writer.WriteString(entity.CurrVersion, 20)

    return writer.Bytes(), nil
}

func (entity *Body_0x12) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.TaskId, err = reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    entity.UpgradeType, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Reason, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.CurrVersion, err = reader.ReadString(20)
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}

func (entity *Body_0x12) MsgID() MsgID {
    return Msg_0x12
}

func (entity *Body_0x11) MsgID() MsgID {
    return Msg_0x11
}

func (entity *Body_0x91) Encode() ([]byte, error) {
    writer := common.NewWriter()

    writer.WriteByte(entity.UpgradeFlag)

    writer.WriteByte(entity.UpgradeType)

    writer.WriteUint32(entity.TaskId)

    writer.WriteString(entity.SrcVersion, 20)

    writer.WriteString(entity.DstVersion, 20)

    writer.WriteByte(entity.UrlLen)

    writer.WriteString(entity.Url)

    writer.WriteString(entity.Md5, 32)

    return writer.Bytes(), nil
}

func (entity *Body_0x91) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.UpgradeFlag, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.UpgradeType, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.TaskId, err = reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    entity.SrcVersion, err = reader.ReadString(20)
    if err != nil {
        return 0, err
    }

    entity.DstVersion, err = reader.ReadString(20)
    if err != nil {
        return 0, err
    }

    entity.UrlLen, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Url, err = reader.ReadString(int(entity.UrlLen))
    if err != nil {
        return 0, err
    }

    entity.Md5, err = reader.ReadString(32)
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
