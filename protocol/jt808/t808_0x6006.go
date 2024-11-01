package jt808

import (
	"github.com/rayjay214/parser/protocol/common"
)

type T808_0x6006 struct {
	AckSeqNo uint16
	Code     uint8
	Content  string
}

func (entity *T808_0x6006) MsgID() MsgID {
	return MsgT808_0x6006
}

func (entity *T808_0x6006) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T808_0x6006) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error

	entity.AckSeqNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Code, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Content, err = reader.ReadString()

	return 0, nil
}
