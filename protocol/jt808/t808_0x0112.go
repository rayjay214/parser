package jt808

import (
	"encoding/json"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/rayjay214/parser/protocol/jt808/errors"
)

// 上传短信修改的省电模式
type T808_0x0112 struct {
	Mode                    byte
	ConnTime                uint16
	IsFlyModeOn             byte
	DurationAfterMonitorCar uint16
}

func (entity T808_0x0112) MarshalJSON() ([]byte, error) {
	type Alias T808_0x0112

	type NewT808_0x0112 struct {
		Alias
		Mode string
	}

	s := NewT808_0x0112{
		Alias: Alias(entity),
	}

	switch entity.Mode {
	case 1:
		s.Mode = "普通模式"
	case 2:
		s.Mode = "周期定位模式"
	case 3:
		s.Mode = "智能省电模式"
	case 4:
		s.Mode = "超级省电模式"
	case 5:
		s.Mode = "智能模式"
	case 6:
		s.Mode = "待机模式"
	case 7:
		s.Mode = "省电模式"
	case 8:
		s.Mode = "点名模式"
	default:
		s.Mode = "其他模式"
	}

	return json.Marshal(s)
}

func (entity *T808_0x0112) MsgID() MsgID {
	return MsgT808_0x0112
}

func (entity *T808_0x0112) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T808_0x0112) Decode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, errors.ErrInvalidBody
	}
	reader := common.NewReader(data)

	var err error

	entity.Mode, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.ConnTime, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.IsFlyModeOn, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.DurationAfterMonitorCar, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return 0, nil
}
