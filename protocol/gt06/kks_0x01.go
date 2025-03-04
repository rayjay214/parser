package gt06

import (
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/common"
	"strconv"
)

type Kks_0x01 struct {
	Proto uint8
	Imei  uint64
	SeqNo uint16
}

func (entity *Kks_0x01) GetSeqNo() uint16 {
	return entity.SeqNo
}

func (entity *Kks_0x01) MsgID() MsgID {
	return Msg_0x01
}

func (entity *Kks_0x01) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteByte(entity.Proto)

	writer.Write(stringToBCD(strconv.FormatUint(entity.Imei, 10), 8))

	writer.WriteUint16(entity.SeqNo)

	return writer.Bytes(), nil
}

func (entity *Kks_0x01) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error
	entity.Proto, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取终端号码
	temp, err := reader.Read(8)
	if err != nil {
		return 0, ErrInvalidHeader
	}
	imei, err := strconv.ParseUint(bcdToString(temp), 10, 64)
	if err != nil {
		return 0, err
	}
	entity.Imei = imei

	entity.SeqNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity Kks_0x01) MarshalJSON() ([]byte, error) {
	type Alias Kks_0x01

	type New0x01 struct {
		Proto string
		Alias
	}

	s := New0x01{
		Alias: Alias(entity),
	}

	s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

	return json.Marshal(s)
}
