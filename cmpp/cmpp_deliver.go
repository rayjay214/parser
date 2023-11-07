package cmpp

import (
    "encoding/hex"
    _ "encoding/json"
    _ "fmt"
    "golang.org/x/text/encoding/unicode"
    "golang.org/x/text/transform"
    "parser/common"
    _ "strconv"
)

type DELIVERREPORT struct {
    MsgSymbol      uint64
    Stat           string
    SubmitTime     string
    DoneTime       string
    DestTerminalId string
    SMMCSequence   uint32
}

type CMPP_DELIVER struct {
    MsgSymbol          uint64
    DestId             string
    ServiceId          string
    TP_pId             uint8
    TP_udhi            uint8
    MsgFmt             uint8
    SrcTerminalId      string
    SrcTerminalType    uint8
    RegisteredDelivery uint8
    MsgLength          uint8
    LongSmsPrefix      string
    MsgContent         string
    DeliverReport      DELIVERREPORT
    LinkID             string
}

func (entity *CMPP_DELIVER) MsgID() MsgID {
    return MSG_CMPP_DELIVER
}

func (entity *CMPP_DELIVER) Encode() ([]byte, error) {
    //todo
    return nil, nil
}

func (entity *CMPP_DELIVER) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.MsgSymbol, err = reader.ReadUint64()
    if err != nil {
        return 0, err
    }

    entity.DestId, err = reader.ReadString(21)
    if err != nil {
        return 0, err
    }

    entity.ServiceId, err = reader.ReadString(10)
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

    entity.SrcTerminalId, err = reader.ReadString(32)
    if err != nil {
        return 0, err
    }

    entity.SrcTerminalType, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.RegisteredDelivery, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.MsgLength, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    if int(entity.RegisteredDelivery) == 1 {
        entity.DeliverReport.MsgSymbol, err = reader.ReadUint64()
        if err != nil {
            return 0, err
        }

        entity.DeliverReport.Stat, err = reader.ReadString(7)
        if err != nil {
            return 0, err
        }

        entity.DeliverReport.SubmitTime, err = reader.ReadString(10)
        if err != nil {
            return 0, err
        }

        entity.DeliverReport.DoneTime, err = reader.ReadString(10)
        if err != nil {
            return 0, err
        }

        entity.DeliverReport.DestTerminalId, err = reader.ReadString(32)
        if err != nil {
            return 0, err
        }

        entity.DeliverReport.SMMCSequence, err = reader.ReadUint32()
        if err != nil {
            return 0, err
        }

    } else {
        var uniCodeContent []byte
        if entity.TP_udhi != 0 {
            var prefix []byte
            prefix, err = reader.Read(6)
            entity.LongSmsPrefix = hex.EncodeToString(prefix)
            if err != nil {
                return 0, err
            }
            uniCodeContent, err = reader.Read(int(entity.MsgLength) - 6)
            if err != nil {
                return 0, err
            }

            utf8Content, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), uniCodeContent)
            entity.MsgContent = string(utf8Content)

        } else {
            uniCodeContent, err = reader.Read(int(entity.MsgLength))
            if err != nil {
                return 0, err
            }

            utf8Content, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), uniCodeContent)
            entity.MsgContent = string(utf8Content)
        }
    }

    entity.LinkID, err = reader.ReadString(20)
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
