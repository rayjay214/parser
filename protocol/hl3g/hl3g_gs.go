package hl3g

import (
	_ "encoding/json"
	_ "fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "github.com/shopspring/decimal"
	_ "strconv"
	_ "time"
)

type HL3G_GS struct {
	Lat  string
	Lng  string
	Time string
}

func (entity *HL3G_GS) MsgID() MsgID {
	return Msg_GS
}

func (entity *HL3G_GS) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteString(entity.Lat)
	writer.WriteString(",")
	writer.WriteString(entity.Lng)
	writer.WriteString(",")
	writer.WriteString(entity.Time)

	return writer.Bytes(), nil
}

func (entity *HL3G_GS) Decode(data []byte) (int, error) {
	//todo
	return len(data), nil
}
