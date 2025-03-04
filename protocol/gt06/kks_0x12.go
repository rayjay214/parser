package gt06

import (
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/shopspring/decimal"
	_ "strconv"
	"time"
)

type Kks_0x12 struct {
	Proto           uint8
	Time            time.Time
	Satellite       uint8 `json:"-"`
	Lat             decimal.Decimal
	Lng             decimal.Decimal
	Speed           uint8
	StatusDirection uint16 `json:"-"`
	Mcc             uint16
	Mnc             uint16
	Lac             uint32 //实际是三个字节
	CellId          uint32
	SeqNo           uint16
}

func (entity *Kks_0x12) GetSeqNo() uint16 {
	return entity.SeqNo
}

func (entity *Kks_0x12) MsgID() MsgID {
	return Msg_0x12
}

func (entity *Kks_0x12) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *Kks_0x12) Decode(data []byte) (int, error) {
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

	entity.Mnc, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	bData := make([]byte, 4)
	bData, _ = reader.Read(3)

	var result uint32
	for _, byteVal := range bData {
		result = (result << 8) | uint32(byteVal)
	}
	entity.Lac = result

	entity.CellId, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	entity.SeqNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity Kks_0x12) MarshalJSON() ([]byte, error) {
	type Alias Kks_0x12

	type New0x12 struct {
		Proto string
		Alias
		Status    map[string]interface{}
		Direction uint16
		GpsNumber uint8
	}

	s := New0x12{
		Alias:  Alias(entity),
		Status: map[string]interface{}{},
	}

	s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

	status := (entity.StatusDirection & 0xFC00) >> 10
	route := entity.StatusDirection & 0x03FF
	gpsNumber := entity.Satellite & 0x0F

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
	s.GpsNumber = gpsNumber

	return json.Marshal(s)
}
