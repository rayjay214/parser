package th

import (
	_ "fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "strconv"
	"strings"
	"time"
)

type TH_LBS struct {
	Proto       string
	Imei        string
	Time        time.Time
	Mcc         string
	Mnc         string
	Lac         string
	Cellid      string
	Rac         string
	Temperature string
	Humility    string
	Status      string
	Battery     string
	GSM         string
	GPSSignal   string
	OnOffStatus string
	Flag        string
}

func (entity *TH_LBS) MsgID() MsgID {
	return Msg_LBS
}

func (entity *TH_LBS) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *TH_LBS) Decode(data []byte) (int, error) {
	strData := string(data)
	strList := strings.Split(strData, "|")
	entity.Proto = strList[0]
	entity.Imei = strList[1]
	strTime := strList[2]
	strLbs := strList[3]
	entity.Temperature = strList[4]
	entity.Humility = strList[5]
	entity.Status = strList[6]
	entity.Battery = strList[7]
	entity.GSM = strList[8]
	entity.GPSSignal = strList[9]
	entity.OnOffStatus = strList[10]
	entity.Flag = strList[11]

	strLbsList := strings.Split(strLbs, ",")
	entity.Mcc = strLbsList[0]
	entity.Mnc = strLbsList[1]
	entity.Cellid = strLbsList[2]
	entity.Lac = strLbsList[3]
	entity.Rac = strLbsList[4]

	entity.Time, _ = time.ParseInLocation("060102150405", strTime, time.UTC)

	return len(data), nil
}
