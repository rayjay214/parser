package ipc

import (
	_ "fmt"
)

type Body_0x8204 struct {
	SeqNo  uint16
	Result uint8
}

func (entity *Body_0x8204) MsgID() MsgID {
	return Msg_0x8204
}

func (entity *Body_0x8204) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteUint16(entity.SeqNo)

	writer.WriteByte(entity.Result)

	return writer.Bytes(), nil
}

func (entity *Body_0x8204) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.SeqNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Result, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
