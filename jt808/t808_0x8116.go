package jt808

import (
	"parser/common"
)

// 终端应答
type T808_0x8116 struct {
	RecordTime byte
}

func (entity *T808_0x8116) MsgID() MsgID {
	return MsgT808_0x8116
}

func (entity *T808_0x8116) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *T808_0x8116) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error

	entity.RecordTime, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
