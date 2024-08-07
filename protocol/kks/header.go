package kks

import (
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "strconv"
)

// 消息头
type Header_0x78 struct {
	Prefix uint16
	MsgLen uint8
}

type Header_0x79 struct {
	Prefix uint16
	MsgLen uint16
}

func (header *Header_0x78) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteUint16(header.Prefix)

	writer.WriteByte(header.MsgLen)

	return writer.Bytes(), nil
}

// 协议解码
func (header *Header_0x78) Decode(data []byte) error {
	if len(data) < 3 {
		return ErrInvalidHeader
	}
	reader := common.NewReader(data)

	prefix, err := reader.ReadUint16()
	if err != nil {
		return ErrInvalidHeader
	}

	msgLen, err := reader.ReadByte()
	if err != nil {
		return ErrInvalidHeader
	}

	header.Prefix = prefix
	header.MsgLen = msgLen

	return nil
}

func (header *Header_0x79) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteUint16(header.Prefix)

	writer.WriteUint16(header.MsgLen)

	return writer.Bytes(), nil
}

// 协议解码
func (header *Header_0x79) Decode(data []byte) error {
	if len(data) < 4 {
		return ErrInvalidHeader
	}
	reader := common.NewReader(data)

	prefix, err := reader.ReadUint16()
	if err != nil {
		return ErrInvalidHeader
	}

	msgLen, err := reader.ReadUint16()
	if err != nil {
		return ErrInvalidHeader
	}

	header.Prefix = prefix
	header.MsgLen = msgLen

	return nil
}

func (header Header_0x79) MarshalJSON() ([]byte, error) {
	type Alias Header_0x79

	return json.Marshal(struct {
		Alias
		Prefix string
	}{
		Alias:  Alias(header),
		Prefix: "0x" + fmt.Sprintf("%04x", header.Prefix),
	})
}

func (header Header_0x78) MarshalJSON() ([]byte, error) {
	type Alias Header_0x78

	return json.Marshal(struct {
		Alias
		Prefix string
	}{
		Alias:  Alias(header),
		Prefix: "0x" + fmt.Sprintf("%04x", header.Prefix),
	})
}
