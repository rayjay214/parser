package ipc

import (
	_ "fmt"
)

type Body_0x1202 struct {
	Appid  string
	Seq    uint16
	Result uint8
}

func (entity *Body_0x1202) MsgID() MsgID {
	return Msg_0x1202
}

func (entity *Body_0x1202) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Appid, 20)

	writer.WriteUint16(entity.Seq)

	writer.WriteByte(entity.Result)

	return writer.Bytes(), nil
}

func (entity *Body_0x1202) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Appid, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.Seq, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Result, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
