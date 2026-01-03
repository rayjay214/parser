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

type HL3G_LK2 struct {
	Step   string
	Roll   string
	Power  string
	Acc    string
	Oil    string
	Charge string
}

func (entity *HL3G_LK2) MsgID() MsgID {
	return Msg_LK2
}

func (entity *HL3G_LK2) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *HL3G_LK2) Decode(data []byte) (int, error) {
	//去掉第一个逗号和最后一个]
	strData := string(data[1 : len(data)-1])
	strList := strings.Split(strData, ",")
	if len(strList) != 6 {
		return 0, ErrInvalidBody
	}
	entity.Step = strList[0]
	entity.Roll = strList[1]
	entity.Power = strList[2]
	entity.Acc = strList[3]
	entity.Oil = strList[4]
	entity.Charge = strList[5]

	return len(data), nil
}
