package tq

import (
	_ "encoding/hex"
	_ "encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/rayjay214/parser/protocol/common"
	"strconv"
	_ "strconv"
	"strings"
	"time"
)

type TQ_Normal struct {
	Prefix    byte
	Imei      string
	Time      time.Time
	Lat       decimal.Decimal
	Lng       decimal.Decimal
	Speed     decimal.Decimal
	Direction string
	Battery   byte
	Status    string
	//Mile uint32
	Mcc    string
	Mnc    string
	Lac    string
	Cellid string
}

func (entity *TQ_Normal) MsgID() MsgID {
	return Msg_Normal
}

func (entity *TQ_Normal) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *TQ_Normal) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error
	entity.Prefix, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	imeiBuf, err := reader.Read(5)
	if err != nil {
		return 0, err
	}
	entity.Imei = fmt.Sprintf("%02x%02x%02x%02x%02x", imeiBuf[0], imeiBuf[1], imeiBuf[2], imeiBuf[3], imeiBuf[4])

	timeBuf, err := reader.Read(3)
	if err != nil {
		return 0, err
	}
	strTime := fmt.Sprintf("%02x%02x%02x", timeBuf[0], timeBuf[1], timeBuf[2])

	dateBuf, err := reader.Read(3)
	if err != nil {
		return 0, err
	}
	strDate := fmt.Sprintf("%02x%02x%02x", dateBuf[0], dateBuf[1], dateBuf[2])
	DD, MM, YY := strDate[0:2], strDate[2:4], strDate[4:6]
	newDate := YY + MM + DD
	newTime := newDate + strTime

	entity.Time, err = time.ParseInLocation("060102150405", newTime, time.UTC)
	if err != nil {
		return 0, err
	}

	fenDivider := decimal.NewFromInt(60)
	latBuf, err := reader.Read(4)
	if err != nil {
		return 0, err
	}
	strLatDu := fmt.Sprintf("%02x", latBuf[0])
	strLatFen := fmt.Sprintf("%02x.%02x%02x", latBuf[1], latBuf[2], latBuf[3])
	latDu, _ := decimal.NewFromString(strLatDu)
	latFen, _ := decimal.NewFromString(strLatFen)
	latFen = latFen.Div(fenDivider)
	entity.Lat = latDu.Add(latFen).Truncate(6)

	entity.Battery, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	if entity.Battery > 240 {
		entity.Battery -= 240
	}

	lngBuf, err := reader.Read(5)
	if err != nil {
		return 0, err
	}
	strLngDu := fmt.Sprintf("%02x%02x", lngBuf[0], lngBuf[1])
	strLngFen := fmt.Sprintf("%02x%02x%02x%02x", (lngBuf[1] & 0x0F), lngBuf[2], lngBuf[3], (lngBuf[4] & 0xF0))
	lngDu, _ := strconv.ParseInt(strLngDu, 10, 64)
	lngFen, _ := strconv.ParseFloat(strLngFen, 64)
	lngDu = lngDu / 10
	lngFen = lngFen / 100000.0 / 60.0
	lng := float64(lngDu) + lngFen
	entity.Lng = decimal.NewFromFloat(lng).Truncate(6)

	speedRouteBuf, err := reader.Read(3)
	if err != nil {
		return 0, err
	}
	strSpeed := fmt.Sprintf("%02x%02x", speedRouteBuf[0], speedRouteBuf[1]&0xF0)
	nmDiv := decimal.NewFromFloat(1.852)
	nmSpeed, _ := decimal.NewFromString(strings.TrimSuffix(strSpeed, "0"))
	entity.Speed = nmSpeed.Mul(nmDiv)

	entity.Direction = fmt.Sprintf("%02x%02x", speedRouteBuf[1]&0x0F, speedRouteBuf[2])

	statusBuf, err := reader.Read(4)
	if err != nil {
		return 0, err
	}
	entity.Status = fmt.Sprintf("%02x%02x%02x%02x", statusBuf[0], statusBuf[1], statusBuf[2], statusBuf[3])

	/*
		entity.Mile, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}
	*/

	_, err = reader.Read(4)
	if err != nil {
		return 0, err
	}

	lbsBuf, err := reader.Read(7)
	if err != nil {
		return 0, err
	}
	entity.Mcc = fmt.Sprintf("%02x%02x", lbsBuf[0], lbsBuf[1])
	entity.Mnc = fmt.Sprintf("%02x", lbsBuf[2])
	entity.Lac = fmt.Sprintf("%02x%02x", lbsBuf[3], lbsBuf[4])
	entity.Cellid = fmt.Sprintf("%02x%02x", lbsBuf[5], lbsBuf[6])

	return len(data), nil
}
