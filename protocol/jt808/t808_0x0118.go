package jt808

import (
	"bytes"
	"github.com/rayjay214/parser/protocol/common"
	"io"
	"io/ioutil"
	"time"
)

// 终端应答
type T808_0x0118 struct {
	PkgSize   byte
	PkgNo     byte
	SessionId uint64
	Time      time.Time
	Packet    io.Reader //录音数据包
}

func (entity *T808_0x0118) MsgID() MsgID {
	return MsgT808_0x0118
}

func (entity *T808_0x0118) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteByte(entity.PkgSize)

	writer.WriteByte(entity.PkgNo)

	writer.WriteUint64(entity.SessionId)

	writer.WriteBcdTime(entity.Time)

	if entity.Packet != nil {
		data, err := ioutil.ReadAll(entity.Packet)
		if err != nil {
			return nil, err
		}
		writer.Write(data)
	}

	return writer.Bytes(), nil
}

func (entity *T808_0x0118) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error

	entity.PkgSize, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.PkgNo, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.SessionId, err = reader.ReadUint64()
	if err != nil {
		return 0, err
	}

	entity.Time, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity *T808_0x0118) DecodePacket(data []byte) error {
	n, err := entity.Decode(data)
	if err != nil {
		return err
	}
	entity.Packet = bytes.NewReader(data[n:])
	return nil
}

func (entity *T808_0x0118) GetTag() uint32 {
	return 0
}

func (entity *T808_0x0118) GetReader() io.Reader {
	return entity.Packet
}

func (entity *T808_0x0118) SetReader(reader io.Reader) {
	entity.Packet = reader
}
