package kks

import (
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "strconv"
	"time"
)

type LbsInfo struct {
	Lac    uint32
	CellId uint64
	Rssi   uint8
}

type Kks_0xa1 struct {
	Proto       uint8
	Time        time.Time
	Mcc         uint16
	Mnc         uint16
	LbsInfoList []LbsInfo
	TimeBefore  uint8
	Lang        uint16
}

func (entity *Kks_0xa1) MsgID() MsgID {
	return Msg_0xa1
}

func (entity *Kks_0xa1) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *Kks_0xa1) Decode(data []byte) (int, error) {
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

	entity.Mnc, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	for i := 0; i < 7; i++ {
		var info LbsInfo
		info.Lac, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}

		info.CellId, err = reader.ReadUint64()
		if err != nil {
			return 0, err
		}

		info.Rssi, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		entity.LbsInfoList = append(entity.LbsInfoList, info)
	}

	return len(data) - reader.Len(), nil
}

func (entity Kks_0xa1) MarshalJSON() ([]byte, error) {
	type Alias Kks_0xa1

	type New0xa1 struct {
		Proto string
		Alias
		Lang string
	}

	s := New0xa1{
		Alias: Alias(entity),
	}

	s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

	if entity.Lang == 0x01 {
		s.Lang = "中文回复"
	} else if entity.Lang == 0x02 {
		s.Lang = "英文回复"
	} else {
		s.Lang = "不需要回复"
	}

	return json.Marshal(s)
}
