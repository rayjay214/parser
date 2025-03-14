package gt06

import (
	"encoding/json"
	"github.com/rayjay214/parser/protocol/common"
	_ "strconv"
	"time"
)

type LbsInfo2G struct {
	Lac    uint16
	CellId uint16
	Rssi   uint8
}

type WifiInfo struct {
	Mac  string
	Rssi uint8
}

type Kks_0x20 struct {
	Proto        uint8
	Time         time.Time
	Mcc          uint16
	Mnc          uint8
	LbsNum       uint8
	LbsInfoList  []LbsInfo2G
	CellId       uint16
	WifiNum      uint8
	WifiInfoList []WifiInfo
	SeqNo        uint16
}

func (entity *Kks_0x20) GetSeqNo() uint16 {
	return entity.SeqNo
}

func (entity *Kks_0x20) MsgID() MsgID {
	return Msg_0x20
}

func (entity *Kks_0x20) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *Kks_0x20) Decode(data []byte) (int, error) {
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

	entity.Mcc, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Mnc, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.LbsNum, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(entity.LbsNum); i++ {
		var info LbsInfo2G
		info.Lac, err = reader.ReadUint16()
		if err != nil {
			return 0, err
		}

		info.CellId, err = reader.ReadUint16()
		if err != nil {
			return 0, err
		}

		info.Rssi, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		entity.LbsInfoList = append(entity.LbsInfoList, info)
	}

	entity.WifiNum, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(entity.WifiNum); i++ {
		var info WifiInfo
		info.Mac, err = reader.ReadString(6)
		if err != nil {
			return 0, err
		}

		info.Rssi, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		entity.WifiInfoList = append(entity.WifiInfoList, info)
	}

	entity.SeqNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity Kks_0x20) MarshalJSON() ([]byte, error) {
	return json.Marshal(entity)
}
