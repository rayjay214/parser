package ipc

import (
	_ "fmt"
)

type Body_0x8021 struct {
	Utc uint32
}

func (entity *Body_0x8021) MsgID() MsgID {
	return Msg_0x8021
}

func (entity *Body_0x8021) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteUint32(entity.Utc)

	return writer.Bytes(), nil
}

func (entity *Body_0x8021) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Utc, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
