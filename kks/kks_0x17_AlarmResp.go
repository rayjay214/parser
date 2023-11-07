package kks

import (
    "encoding/json"
    "fmt"
    "golang.org/x/text/encoding/unicode"
    "golang.org/x/text/transform"
    "parser/common"
    _ "strconv"
)

type Kks_0x17_AlarmResp struct {
    Proto        uint8
    ContentLen   uint8
    AlarmFlag    uint32
    AlarmContent string
    Phone        string
}

func (entity *Kks_0x17_AlarmResp) MsgID() MsgID {
    return Msg_0x17_AlarmResp
}

func (entity *Kks_0x17_AlarmResp) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *Kks_0x17_AlarmResp) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.Proto, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.ContentLen, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.AlarmFlag, err = reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    _, err = reader.ReadString(8) //ADDRESS
    if err != nil {
        return 0, err
    }

    _, err = reader.ReadString(2) //&&
    if err != nil {
        return 0, err
    }

    addrLen := int(entity.ContentLen) - 4 - 8 - 2 - 2 - 21 - 2
    var uniCodeContent []byte
    uniCodeContent, err = reader.Read(int(addrLen))
    if err != nil {
        return 0, err
    }

    utf8Content, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), uniCodeContent)
    entity.AlarmContent = string(utf8Content)

    _, err = reader.ReadString(2) //&&
    if err != nil {
        return 0, err
    }

    entity.Phone, err = reader.ReadString(21)
    if err != nil {
        return 0, err
    }

    _, err = reader.ReadString(2) //##
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}

func (entity Kks_0x17_AlarmResp) MarshalJSON() ([]byte, error) {
    type Alias Kks_0x17_AlarmResp

    type New0x17_AlarmResp struct {
        Proto     string
        AlarmFlag string
        Alias
    }

    s := New0x17_AlarmResp{
        Alias: Alias(entity),
    }

    s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)
    s.AlarmFlag = "0x" + fmt.Sprintf("%02x", entity.AlarmFlag)

    return json.Marshal(s)
}
