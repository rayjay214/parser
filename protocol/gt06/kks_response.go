package gt06

import (
	"github.com/rayjay214/parser/protocol/common"
	_ "strconv"
)

type KksResponse struct {
	Proto uint8
	SeqNo uint16
}

func (entity *KksResponse) MsgID() MsgID {
	return Msg_0x01
}

func (entity *KksResponse) GetSeqNo() uint16 {
	return entity.SeqNo
}

func (entity *KksResponse) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteByte(entity.Proto)

	writer.WriteUint16(entity.SeqNo)

	return writer.Bytes(), nil
}

func (entity *KksResponse) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error
	entity.Proto, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.SeqNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
