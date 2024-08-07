package tq

import (
	"encoding/hex"
	"encoding/json"
	_ "fmt"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/shopspring/decimal"
	_ "strconv"
	"strings"
	"time"
)

type TQ_V2 struct {
	Manu      string
	Imei      string
	Proto     string
	Time      time.Time
	Located   string
	Lat       decimal.Decimal
	Lng       decimal.Decimal
	Speed     decimal.Decimal
	Direction string
	Status    string
	Mcc       string
	Mnc       string
	Lac       string
	Cellid    string
}

func (entity *TQ_V2) MsgID() MsgID {
	return Msg_V2
}

func (entity *TQ_V2) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *TQ_V2) Decode(data []byte) (int, error) {
	strData := string(data)
	strList := strings.Split(strData, ",")
	entity.Manu = strings.TrimPrefix(strList[0], "*")
	entity.Imei = strList[1]
	entity.Proto = strList[2]
	strTime := strList[3]
	entity.Located = strList[4]
	strLat := strList[5]
	strLatFlag := strList[6]
	strLng := strList[7]
	strLngFlag := strList[8]
	strSpeed := strList[9]
	entity.Direction = strList[10]
	strDate := strList[11]
	entity.Status = strList[12]
	entity.Mcc = strList[13]
	entity.Mnc = strList[14]
	entity.Lac = strList[15]
	entity.Cellid = strings.TrimSuffix(strList[16], "#")

	//transform lat lng
	fenDivider := decimal.NewFromInt(60)
	strlatDu := strLat[0:2]
	strlatFen := strLat[2:]
	latDu, _ := decimal.NewFromString(strlatDu)
	latFen, _ := decimal.NewFromString(strlatFen)
	latFen = latFen.Div(fenDivider)
	entity.Lat = latDu.Add(latFen).Truncate(6)
	if strLatFlag == "S" {
		entity.Lat = decimal.Zero.Sub(entity.Lat)
	}

	strlngDu := strLng[0:3]
	strlngFen := strLng[3:]
	lngDu, _ := decimal.NewFromString(strlngDu)
	lngFen, _ := decimal.NewFromString(strlngFen)
	lngFen = lngFen.Div(fenDivider)
	entity.Lng = lngDu.Add(lngFen).Truncate(6)
	if strLngFlag == "W" {
		entity.Lng = decimal.Zero.Sub(entity.Lng)
	}

	//transform speed
	nmDiv := decimal.NewFromFloat(1.852)
	nmSpeed, _ := decimal.NewFromString(strSpeed)
	entity.Speed = nmSpeed.Mul(nmDiv)

	//parse time
	DD, MM, YY := strDate[0:2], strDate[2:4], strDate[4:6]
	newDate := YY + MM + DD
	newTime := newDate + strTime
	entity.Time, _ = time.ParseInLocation("060102150405", newTime, time.UTC)

	return len(data), nil
}

func (entity TQ_V2) MarshalJSON() ([]byte, error) {
	type Alias TQ_V2

	type NewV2 struct {
		Alias
		Status []string
	}

	s := NewV2{
		Alias: Alias(entity),
	}

	s.Status = make([]string, 0)

	fGetbit := func(value []byte, offset int) uint32 {
		return uint32(value[0]) & (1 << offset) >> offset
	}

	first, second, third, fouth := entity.Status[0:2], entity.Status[2:4], entity.Status[4:6], entity.Status[6:8]
	bFirst, _ := hex.DecodeString(first)
	bSecond, _ := hex.DecodeString(second)
	bThird, _ := hex.DecodeString(third)
	bFouth, _ := hex.DecodeString(fouth)

	if fGetbit(bFirst, 1) == 0 {
		s.Status = append(s.Status, "位移报警")
	}
	if fGetbit(bFirst, 2) == 0 {
		s.Status = append(s.Status, "补报数据")
	}
	if fGetbit(bFirst, 3) == 0 {
		s.Status = append(s.Status, "断油电状态")
	}
	if fGetbit(bFirst, 4) == 0 {
		s.Status = append(s.Status, "电瓶拆除报警")
	}
	if fGetbit(bSecond, 1) == 0 {
		s.Status = append(s.Status, "震动报警")
	}
	if fGetbit(bSecond, 3) == 0 {
		s.Status = append(s.Status, "主机掉电由后备电池供电")
	}
	if fGetbit(bThird, 0) == 0 {
		s.Status = append(s.Status, "车门开")
	}
	if fGetbit(bThird, 1) == 0 {
		s.Status = append(s.Status, "车辆设防")
	}
	if fGetbit(bThird, 2) == 0 {
		s.Status = append(s.Status, "ACC关")
	}
	if fGetbit(bFouth, 1) == 0 {
		s.Status = append(s.Status, "劫警")
	}
	if fGetbit(bFouth, 2) == 0 {
		s.Status = append(s.Status, "超速报警")
	}
	if fGetbit(bFouth, 3) == 0 {
		s.Status = append(s.Status, "非法点火报警")
	}
	if fGetbit(bFouth, 4) == 0 {
		s.Status = append(s.Status, "非法开门报警")
	}
	if fGetbit(bFouth, 5) == 0 {
		s.Status = append(s.Status, "电瓶低电报警")
	}
	if fGetbit(bFouth, 6) == 0 {
		s.Status = append(s.Status, "电池低电报警")
	}

	return json.Marshal(s)
}
