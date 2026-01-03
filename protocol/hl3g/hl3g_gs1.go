package hl3g

import (
	_ "encoding/json"
	_ "fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "github.com/shopspring/decimal"
	_ "strconv"
	"strings"
	_ "time"
)

type HL3G_GS1 struct {
	Reserve1 string
	Reserve2 string
	Mcc      string
	Mnc      string
	Lac      string
	CellId   string
	Rssi     string
}

func (entity *HL3G_GS1) MsgID() MsgID {
	return Msg_GS1
}

func (entity *HL3G_GS1) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *HL3G_GS1) Decode(data []byte) (int, error) {
	//去掉第一个逗号和最后一个]
	strData := string(data[1 : len(data)-1])
	strList := strings.Split(strData, ",")
	if len(strList) != 7 {
		return 0, ErrInvalidBody
	}
	entity.Reserve1 = strList[0]
	entity.Reserve2 = strList[1]
	entity.Mcc = strList[2]
	entity.Mnc = strList[3]
	entity.Lac = strList[4]
	entity.CellId = strList[5]
	entity.Rssi = strList[6]

	return len(data), nil
}
