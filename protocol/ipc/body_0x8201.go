package ipc

import (
	_ "fmt"
)

type Body_0x8201 struct {
	Appid     string
	BeginTime uint32
	EndTime   uint32
}

func (entity *Body_0x8201) MsgID() MsgID {
	return Msg_0x8201
}

func (entity *Body_0x8201) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Appid, 20)

	writer.WriteUint32(entity.BeginTime)

	writer.WriteUint32(entity.EndTime)

	return writer.Bytes(), nil
}

func (entity *Body_0x8201) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Appid, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.BeginTime, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	entity.EndTime, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
