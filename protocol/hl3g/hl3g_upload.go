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

type HL3G_UPLOAD struct {
	Interval string
}

func (entity *HL3G_UPLOAD) MsgID() MsgID {
	return Msg_UPLOAD
}

func (entity *HL3G_UPLOAD) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteString(entity.Interval)

	return writer.Bytes(), nil
}

func (entity *HL3G_UPLOAD) Decode(data []byte) (int, error) {
	//去掉第一个逗号和最后一个]
	strData := string(data[1 : len(data)-1])
	strList := strings.Split(strData, ",")
	entity.Interval = strList[0]

	return len(data), nil
}
