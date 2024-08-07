package jt808

import (
	"github.com/rayjay214/parser/protocol/common"
)

// 平台通用应答
type T808_0x0818 struct {
	PhoneLen byte
	Phone    string
}

func (entity *T808_0x0818) MsgID() MsgID {
	return MsgT808_0x0818
}

func (entity *T808_0x0818) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo
	return writer.Bytes(), nil
}

func (entity *T808_0x0818) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error
	entity.PhoneLen, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Phone, err = reader.ReadString(int(entity.PhoneLen))
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}