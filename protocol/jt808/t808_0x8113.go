package jt808

import (
	"github.com/rayjay214/parser/protocol/common"
)

type T808_0x8113 struct {
	PhoneLen   byte
	Phone      string
	SmsLen     byte
	SmsContent string
}

func (entity *T808_0x8113) MsgID() MsgID {
	return MsgT808_0x8113
}

func (entity *T808_0x8113) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T808_0x8113) Decode(data []byte) (int, error) {
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

	entity.SmsLen, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	SmsContent, err := reader.Read(int(entity.SmsLen))
	if err != nil {
		return 0, err
	}
	entity.SmsContent = common.BytesToString(SmsContent[:])

	return 0, nil
}
