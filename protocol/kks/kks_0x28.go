package kks

import (
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "strconv"
	"time"
)

type Kks_0x28 struct {
	Proto      uint8
	Time       time.Time
	Mcc        uint16
	Mnc        uint8
	Lac        []uint16
	CellId     []uint16
	Rssi       []int
	TimeBefore uint8
	Lang       uint16
}

func (entity *Kks_0x28) MsgID() MsgID {
	return Msg_0x28
}

func (entity *Kks_0x28) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *Kks_0x28) Decode(data []byte) (int, error) {
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

	entity.Lac = make([]uint16, 0)
	entity.CellId = make([]uint16, 0)
	entity.Rssi = make([]int, 0)
	//固定上传7个
	var lac, cellid uint16
	var rssi uint8
	for i := 0; i < 7; i++ {
		lac, err = reader.ReadUint16()
		if err != nil {
			return 0, err
		}
		entity.Lac = append(entity.Lac, lac)

		reader.ReadByte() //useless

		cellid, err = reader.ReadUint16()
		if err != nil {
			return 0, err
		}
		entity.CellId = append(entity.CellId, cellid)

		rssi, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}
		entity.Rssi = append(entity.Rssi, int(rssi))
	}

	entity.TimeBefore, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Lang, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity Kks_0x28) MarshalJSON() ([]byte, error) {
	type Alias Kks_0x28

	type New0x28 struct {
		Proto string
		Lang  string
		Alias
	}

	s := New0x28{
		Alias: Alias(entity),
	}

	s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

	if GetBit(int(entity.Lang), 0) == 1 {
		s.Lang = "中文"
	} else {
		s.Lang = "英文"
	}

	return json.Marshal(s)
}
