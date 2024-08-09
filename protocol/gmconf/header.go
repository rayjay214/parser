package gmconf

import (
	"strconv"
	//"fmt"
	"github.com/rayjay214/parser/protocol/common"
)

// 消息头
type Header struct {
	Prefix uint16
	MsgId  MsgID
	MsgLen uint16
	Imei   uint64
}

// 协议编码
func (header *Header) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteUint16(uint16(0x6666))

	writer.WriteByte(uint8(header.MsgId))

	writer.WriteUint16(header.MsgLen)

	writer.Write(common.StringToBCD(strconv.FormatUint(header.Imei, 10), 8))

	return writer.Bytes(), nil
}

// 协议解码
func (header *Header) Decode(data []byte) error {
	if len(data) < MessageHeaderSize {
		return ErrInvalidHeader
	}
	reader := common.NewReader(data)

	prefix, err := reader.ReadUint16()
	if err != nil {
		return ErrInvalidHeader
	}

	msgId, err := reader.ReadByte()
	if err != nil {
		return ErrInvalidHeader
	}

	msgLen, err := reader.ReadUint16()
	if err != nil {
		return ErrInvalidHeader
	}

	temp, err := reader.Read(8)
	if err != nil {
		return ErrInvalidHeader
	}
	imei, err := strconv.ParseUint(bcdToString(temp), 10, 64)
	if err != nil {
		return err
	}

	header.Prefix = prefix
	header.MsgId = MsgID(msgId)
	header.MsgLen = msgLen
	header.Imei = imei

	return nil
}
