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

type HL3G_VERNO struct {
	Version string
}

func (entity *HL3G_VERNO) MsgID() MsgID {
	return Msg_VERNO
}

func (entity *HL3G_VERNO) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteString(entity.Version)

	return writer.Bytes(), nil
}

func (entity *HL3G_VERNO) Decode(data []byte) (int, error) {
	//去掉第一个逗号和最后一个]
	if len(data) > 1 {
		strData := string(data[1 : len(data)-1])
		strList := strings.Split(strData, ",")
		entity.Version = strList[0]
	}

	return len(data), nil
}
