package kks

import (
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/shopspring/decimal"
	_ "strconv"
	"time"
)

type Kks_0x22 struct {
	Proto           uint8
	Satellite       uint8 `json:"-"`
	Lat             decimal.Decimal
	Lng             decimal.Decimal
	Speed           uint8
	StatusDirection uint16 `json:"-"`
	Time            time.Time
	Mcc             uint16
	Mnc             uint8
	Lac             uint16
	CellId          uint16
	Acc             uint8
	Mode            uint8
	GpsMode         uint8
}

func (entity *Kks_0x22) MsgID() MsgID {
	return Msg_0x22
}

func (entity *Kks_0x22) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *Kks_0x22) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error
	entity.Proto, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Time, err = reader.ReadStrTime()
	if err != nil {
		return 0, err
	}

	entity.Satellite, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	latitude, err := reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	longitude, err := reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	entity.Speed, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	statusRoute, err := reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	status := (statusRoute & 0xFC00) >> 10
	//route := statusRoute & 0x03FF

	entity.StatusDirection = statusRoute

	var south, west bool
	if GetBit(int(status), 0) == 0 {
		south = true
	}
	if GetBit(int(status), 1) == 1 {
		west = true
	}

	entity.Lat, entity.Lng = getGeoPoint(latitude, south, longitude, west)

	entity.Mcc, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Mnc, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Lac, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	reader.ReadByte()

	entity.CellId, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Acc, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Mode, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.GpsMode, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity Kks_0x22) MarshalJSON() ([]byte, error) {
	type Alias Kks_0x22

	type New0x22 struct {
		Proto string
		Alias
		Status    map[string]interface{}
		Direction uint16
		Acc       string
		GpsMode   string
	}

	s := New0x22{
		Alias:  Alias(entity),
		Status: map[string]interface{}{},
	}

	s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

	status := (entity.StatusDirection & 0xFC00) >> 10
	route := entity.StatusDirection & 0x03FF

	statusMap := map[string]interface{}{}

	if GetBit(int(status), 2) == 1 {
		statusMap["定位状态"] = "GPS已定位"
	} else {
		statusMap["定位状态"] = "GPS未定位"
	}
	if GetBit(int(status), 3) == 1 {
		statusMap["定位类型"] = "差分定位"
	} else {
		statusMap["定位类型"] = "GPS实时"
	}
	if GetBit(int(status), 4) == 1 {
		statusMap["设备状态"] = "运动"
	} else {
		statusMap["设备状态"] = "静止"
	}
	s.Status = statusMap
	s.Direction = route

	if entity.Acc == 1 {
		s.Acc = "高"
	} else {
		s.Acc = "低"
	}

	if entity.GpsMode == 1 {
		s.GpsMode = "补传"
	} else {
		s.GpsMode = "实时上传"
	}

	return json.Marshal(s)
}
