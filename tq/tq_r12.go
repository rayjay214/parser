package tq

import (
    _ "encoding/json"
    _ "fmt"
    _ "github.com/shopspring/decimal"
    "parser/common"
    _ "strconv"
    "strings"
    _ "time"
)

type TQ_R12 struct {
    Manu  string
    Imei  string
    Proto string
    Time  string
}

func (entity *TQ_R12) MsgID() MsgID {
    return Msg_R12
}

func (entity *TQ_R12) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *TQ_R12) Decode(data []byte) (int, error) {
    strData := string(data)
    strList := strings.Split(strData, ",")
    entity.Manu = strings.TrimPrefix(strList[0], "*")
    entity.Imei = strList[1]
    entity.Proto = strList[2]
    entity.Time = strings.TrimSuffix(strList[3], "#")

    return len(data), nil
}
