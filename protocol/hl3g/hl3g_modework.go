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

type HL3G_MODEWORK struct {
	Mode string
}

func (entity *HL3G_MODEWORK) MsgID() MsgID {
	return Msg_MODEWORK
}

func (entity *HL3G_MODEWORK) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteString(entity.Mode)

	return writer.Bytes(), nil
}

func (entity *HL3G_MODEWORK) Decode(data []byte) (int, error) {
	//去掉第一个逗号和最后一个]
	if len(data) > 1 {
		strData := string(data[1 : len(data)-1])
		strList := strings.Split(strData, ",")
		entity.Mode = strList[0]
	}

	return len(data), nil
}
