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

type HL3G_STU struct {
	WorkPattern     string
	Sos             string
	Record          string
	Shake           string
	PersistRecord   string
	AirplaneMode    string
	RecordDuration  string
	InflectionPoint string
	Dismantle       string
	OnOff           string
}

func (entity *HL3G_STU) MsgID() MsgID {
	return Msg_STU
}

func (entity *HL3G_STU) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *HL3G_STU) Decode(data []byte) (int, error) {
	strData := string(data[1 : len(data)-1])
	strList := strings.Split(strData, ";")
	entity.WorkPattern = strList[0]
	entity.Sos = strList[1]
	entity.Record = strList[2]
	entity.Shake = strList[3]
	entity.PersistRecord = strList[4]
	entity.AirplaneMode = strList[5]
	entity.RecordDuration = strList[6]
	entity.InflectionPoint = strList[7]
	entity.Dismantle = strList[8]
	entity.OnOff = strList[9]

	return len(data), nil
}
