package jt808

import (
	"github.com/rayjay214/parser/protocol/common"
	"time"
)

// 请求同步时间
type T808_0x1005 struct {
	Voltage uint16
	Current uint16
	Reserve uint16
	Time    time.Time
}

func (entity *T808_0x1005) MsgID() MsgID {
	return MsgT808_0x1005
}

func (entity *T808_0x1005) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T808_0x1005) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error

	entity.Voltage, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Current, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Reserve, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Time, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}

	return 0, nil
}
