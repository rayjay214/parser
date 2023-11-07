package cmpp

import (
    _ "encoding/json"
    _ "fmt"
    "github.com/rayjay214/parser/common"
    _ "strconv"
)

type CMPP_SUBMIT struct {
    MsgSymbol          uint64
    PkTotal            uint8
    PKNumber           uint8
    RegisteredDelivery uint8
    MsgLevel           uint8
    ServiceId          string
    FeeUserType        uint8
    FeeTerminalId      string
    FeeTerminalType    uint8
    TP_pId             uint8
    TP_udhi            uint8
    MsgFmt             uint8
    MsgSrc             string
    FeeType            string
    FeeCode            string
    ValidTime          string
    AtTime             string
    SrcId              string
    DestUsrTl          uint8
    DestTerminalId     string
    DestTerminalType   uint8
    MsgLength          uint8
    MsgContent         string
    LinkID             string
}

func (entity *CMPP_SUBMIT) MsgID() MsgID {
    return MSG_CMPP_SUBMIT
}

func (entity *CMPP_SUBMIT) Encode() ([]byte, error) {
    //todo
    return nil, nil
}

func (entity *CMPP_SUBMIT) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.MsgSymbol, err = reader.ReadUint64()
    if err != nil {
        return 0, err
    }

    entity.PkTotal, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.PKNumber, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.RegisteredDelivery, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.MsgLevel, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.ServiceId, err = reader.ReadString(10)
    if err != nil {
        return 0, err
    }

    entity.FeeUserType, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.FeeTerminalId, err = reader.ReadString(32)
    if err != nil {
        return 0, err
    }

    entity.FeeTerminalType, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.TP_pId, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.TP_udhi, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.MsgFmt, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.MsgSrc, err = reader.ReadString(6)
    if err != nil {
        return 0, err
    }

    entity.FeeType, err = reader.ReadString(2)
    if err != nil {
        return 0, err
    }

    entity.FeeCode, err = reader.ReadString(6)
    if err != nil {
        return 0, err
    }

    entity.ValidTime, err = reader.ReadString(17)
    if err != nil {
        return 0, err
    }

    entity.AtTime, err = reader.ReadString(17)
    if err != nil {
        return 0, err
    }

    entity.SrcId, err = reader.ReadString(21)
    if err != nil {
        return 0, err
    }

    entity.DestUsrTl, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.DestTerminalId, err = reader.ReadString(32 * int(entity.DestUsrTl))
    if err != nil {
        return 0, err
    }

    entity.DestTerminalType, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.MsgLength, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.MsgContent, err = reader.ReadString(int(entity.MsgLength))
    if err != nil {
        return 0, err
    }

    entity.LinkID, err = reader.ReadString(20)
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
