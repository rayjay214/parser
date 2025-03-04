package gt06

import (
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "strconv"
)

// 消息头
type Header struct {
	Prefix uint16
	MsgLen interface{}
}

func (header *Header) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteUint16(header.Prefix)

	if msgLen, ok := header.MsgLen.(byte); ok {
		writer.WriteByte(msgLen)
	}

	if msgLen, ok := header.MsgLen.(uint16); ok {
		writer.WriteUint16(msgLen)
	}

	return writer.Bytes(), nil
}

// 协议解码
func (header *Header) Decode(data []byte) error {
	if len(data) < 2 {
		return ErrInvalidHeader
	}
	reader := common.NewReader(data)

	prefix, err := reader.ReadUint16()
	if err != nil {
		return ErrInvalidHeader
	}

	var msgLen interface{}
	if prefix == 0x7979 {
		msgLen, err = reader.ReadUint16()
		if err != nil {
			return ErrInvalidHeader
		}
	} else {
		msgLen, err = reader.ReadByte()
		if err != nil {
			return ErrInvalidHeader
		}
	}

	header.Prefix = prefix
	header.MsgLen = msgLen

	return nil
}

func (header Header) MarshalJSON() ([]byte, error) {
	type Alias Header

	return json.Marshal(struct {
		Alias
		Prefix string
	}{
		Alias:  Alias(header),
		Prefix: "0x" + fmt.Sprintf("%04x", header.Prefix),
	})
}
