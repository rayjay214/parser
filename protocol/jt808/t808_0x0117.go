package jt808

import (
	"bytes"
	"github.com/rayjay214/parser/protocol/common"
	"io"
	"io/ioutil"
)

// 终端应答
type T808_0x0117 struct {
	PkgSize   byte //总包个数
	PkgNo     byte //包序号
	SessionId string
	Packet    io.Reader //录音数据包
}

func (entity *T808_0x0117) MsgID() MsgID {
	return MsgT808_0x0117
}

func (entity *T808_0x0117) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo
	writer.WriteByte(entity.PkgSize)

	writer.WriteByte(entity.PkgNo)

	writer.WriteString(entity.SessionId)

	// 写入数据包
	if entity.Packet != nil {
		data, err := ioutil.ReadAll(entity.Packet)
		if err != nil {
			return nil, err
		}
		writer.Write(data)
		//log.Infof("encode insert pkt %v", len(data))
	}

	return writer.Bytes(), nil
}

func (entity *T808_0x0117) Decode(data []byte) (int, error) {
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

	entity.SessionId, err = reader.ReadString(8)
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity *T808_0x0117) DecodePacket(data []byte) error {
	n, err := entity.Decode(data)
	if err != nil {
		return err
	}
	//log.Infof("decode insert pkt %v", len(data)-n)
	entity.Packet = bytes.NewReader(data[n:])
	return nil
}

func (entity *T808_0x0117) GetTag() uint32 {
	return 0
}

func (entity *T808_0x0117) GetReader() io.Reader {
	return entity.Packet
}

func (entity *T808_0x0117) SetReader(reader io.Reader) {
	entity.Packet = reader
}
