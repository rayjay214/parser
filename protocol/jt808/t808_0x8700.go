package jt808

import (
	"github.com/rayjay214/parser/protocol/common"
	"github.com/rayjay214/parser/protocol/jt808/errors"
)

// 行驶记录数据采集命令
type T808_0x8700 struct {
	// 命令字
	Cmd byte
	// 数据块
	Data []byte
}

func (entity *T808_0x8700) MsgID() MsgID {
	return MsgT808_0x8700
}

func (entity *T808_0x8700) Encode() ([]byte, error) {
	writer := common.NewWriter()

	// 写入命令字
	writer.WriteByte(entity.Cmd)

	// 写入数据块
	writer.Write(entity.Data)
	return writer.Bytes(), nil
}

func (entity *T808_0x8700) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, errors.ErrInvalidBody
	}
	reader := common.NewReader(data)

	// 读取命令字
	var err error
	entity.Cmd, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取数据块
	entity.Data, err = reader.Read()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
