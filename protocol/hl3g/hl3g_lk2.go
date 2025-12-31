package hl3g

import (
	"encoding/hex"
	_ "encoding/json"
	_ "fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "github.com/shopspring/decimal"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
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
