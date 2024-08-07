package jt808

import (
	"github.com/rayjay214/parser/protocol/common"
	"github.com/rayjay214/parser/protocol/jt808/errors"
)

// 终端应答
type T808_0x8109 struct {
	Year   uint16
	Month  byte
	Day    byte
	Hour   byte
	Minute byte
	Second byte
	Result byte
}

func (entity *T808_0x8109) MsgID() MsgID {
	return MsgT808_0x8109
}

func (entity *T808_0x8109) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo
	writer.WriteUint16(entity.Year)

	writer.WriteByte(entity.Month)

	writer.WriteByte(entity.Day)

	writer.WriteByte(entity.Hour)

	writer.WriteByte(entity.Minute)

	writer.WriteByte(entity.Second)

	writer.WriteByte(entity.Result)

	return writer.Bytes(), nil
}

func (entity *T808_0x8109) Decode(data []byte) (int, error) {
	if len(data) < 8 {
		return 0, errors.ErrInvalidBody
	}
	reader := common.NewReader(data)

	var err error

	entity.Year, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Month, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Day, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Hour, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Minute, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Second, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	/*
		entity.Result, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}
	*/

	return len(data) - reader.Len(), nil
}
