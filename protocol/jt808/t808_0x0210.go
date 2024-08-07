package jt808

import (
	"github.com/rayjay214/parser/protocol/common"
	"time"
)

// 休眠电池电量更新
type T808_0x0210 struct {
	Battery byte
	Time    time.Time
}

func (entity *T808_0x0210) MsgID() MsgID {
	return MsgT808_0x0210
}

func (entity *T808_0x0210) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *T808_0x0210) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error

	entity.Battery, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Time, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}

	return 0, nil
}
