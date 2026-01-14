package hl3g

import (
	_ "encoding/json"
	_ "fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "github.com/shopspring/decimal"
	"strconv"
	_ "strconv"
	"strings"
	_ "time"
)

type HL3G_UD2 struct {
	LocInfo LocationInfo
}

func (entity *HL3G_UD2) MsgID() MsgID {
	return Msg_UD2
}

func (entity *HL3G_UD2) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *HL3G_UD2) Decode(data []byte) (int, error) {
	//去掉第一个逗号和最后一个]
	strData := string(data[1 : len(data)-1])
	strList := strings.Split(strData, ",")
	entity.LocInfo.Date = strList[0]
	entity.LocInfo.Time = strList[1]
	entity.LocInfo.Located = strList[2]
	entity.LocInfo.Lat = strList[3]
	entity.LocInfo.LatFlag = strList[4]
	entity.LocInfo.Lng = strList[5]
	entity.LocInfo.LngFlag = strList[6]
	entity.LocInfo.Speed = strList[7]
	entity.LocInfo.Direction = strList[8]
	entity.LocInfo.Altitude = strList[9]
	entity.LocInfo.StarNum = strList[10]
	entity.LocInfo.Gsm = strList[11]
	entity.LocInfo.Power = strList[12]
	entity.LocInfo.Step = strList[13]
	entity.LocInfo.Roll = strList[14]
	entity.LocInfo.Status = strList[15]
	entity.LocInfo.LbsNum = strList[16]

	lbsCount, _ := strconv.Atoi(entity.LocInfo.LbsNum)
	for i := 0; i < lbsCount; i++ {
		var lbs LbsInfo
		lbs.Ta = strList[17]
		lbs.Mcc = strList[18]
		lbs.Mnc = strList[19]
		lbs.Lac = strList[20+i*3]
		lbs.CellId = strList[21+i*3]
		lbs.Rssi = strList[22+i*3]
		entity.LocInfo.Lbs = append(entity.LocInfo.Lbs, lbs)
	}
	var currentIndex int
	if lbsCount >= 0 {
		currentIndex = 17 + 3 + lbsCount*3
	} else {
		currentIndex = 17
	}

	entity.LocInfo.WifiNum = strList[currentIndex]
	wifiCount, _ := strconv.Atoi(entity.LocInfo.WifiNum)
	for i := 0; i < wifiCount; i++ {
		var wifi WifiInfo
		wifi.Name = strList[currentIndex+1+i*3]
		wifi.Mac = strList[currentIndex+2+i*3]
		wifi.Rssi = strList[currentIndex+3+i*3]
		entity.LocInfo.Wifi = append(entity.LocInfo.Wifi, wifi)
	}
	//entity.LocInfo.LocAccuracy = strList[currentIndex+1+wifiCount*3]

	return len(data), nil
}
