package tq

import (
    "encoding/hex"
    _ "encoding/json"
    _ "fmt"
    _ "github.com/shopspring/decimal"
    "golang.org/x/text/encoding/unicode"
    "golang.org/x/text/transform"
    "github.com/rayjay214/parser/common"
    _ "strconv"
    "strings"
    _ "time"
)

type TQ_I1 struct {
    Manu        string
    Imei        string
    Proto       string
    Time        string
    DisplayTime string
    Code        string
    InfoLen     string
    Info        string
}

func (entity *TQ_I1) MsgID() MsgID {
    return Msg_V1
}

func (entity *TQ_I1) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *TQ_I1) Decode(data []byte) (int, error) {
    strData := string(data)
    strList := strings.Split(strData, ",")
    entity.Manu = strings.TrimPrefix(strList[0], "*")
    entity.Imei = strList[1]
    entity.Proto = strList[2]
    entity.Time = strList[3]
    entity.DisplayTime = strList[4]
    entity.Code = strList[5]
    entity.InfoLen = strList[6]
    uniCodeInfo := strings.TrimSuffix(strList[7], "#")
    bUniCodeInfo, _ := hex.DecodeString(uniCodeInfo)

    utf8Content, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), bUniCodeInfo)
    entity.Info = string(utf8Content)

    return len(data), nil
}
