package jt808

import (
	"github.com/rayjay214/parser/protocol/common"
	"time"
)

// 终端应答
type T808_0x8118 struct {
	PkgNo     byte
	SessionId uint64
	Time      time.Time
}

func (entity *T808_0x8118) MsgID() MsgID {
	return MsgT808_0x8118
}

func (entity *T808_0x8118) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteByte(entity.PkgNo)

	writer.WriteUint64(entity.SessionId)

	writer.WriteBcdTime(entity.Time)

	return writer.Bytes(), nil
}

func (entity *T808_0x8118) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error

	entity.PkgNo, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	_, err = reader.ReadUint64()
	if err != nil {
		return 0, err
	}

	entity.Time, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
