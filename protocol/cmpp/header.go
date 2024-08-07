package cmpp

import (
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "strconv"
)

// 消息头
type Header struct {
	TotalLen uint32
	CmdId    MsgID
	SeqId    uint32
}

// 协议编码
func (header *Header) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteUint32(header.TotalLen)

	writer.WriteUint32(uint32(header.CmdId))

	writer.WriteUint32(header.SeqId)

	return writer.Bytes(), nil
}

// 协议解码
func (header *Header) Decode(data []byte) error {
	if len(data) < 12 {
		return ErrInvalidHeader
	}
	reader := common.NewReader(data)

	msgLen, err := reader.ReadUint32()
	if err != nil {
		return ErrInvalidHeader
	}

	cmdId, err := reader.ReadUint32()
	if err != nil {
		return ErrInvalidHeader
	}

	seqId, err := reader.ReadUint32()
	if err != nil {
		return ErrInvalidHeader
	}

	header.TotalLen = msgLen
	header.CmdId = MsgID(cmdId)
	header.SeqId = seqId

	return nil
}

func (header Header) MarshalJSON() ([]byte, error) {
	type Alias Header

	return json.Marshal(struct {
		Alias
		CmdId string
	}{
		Alias: Alias(header),
		CmdId: "0x" + fmt.Sprintf("%08x", header.CmdId),
	})
}
