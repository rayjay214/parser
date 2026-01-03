package hl3g

import (
	_ "encoding/json"
	_ "fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "github.com/shopspring/decimal"
	_ "strconv"
	_ "time"
)

type HL3G_CCID struct {
	Iccid string
}

func (entity *HL3G_CCID) MsgID() MsgID {
	return Msg_CCID
}

func (entity *HL3G_CCID) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *HL3G_CCID) Decode(data []byte) (int, error) {
	//去掉第一个逗号和最后一个]
	strData := string(data[1 : len(data)-1])
	entity.Iccid = strData

	return len(data), nil
}
