package ipc

import (
	_ "fmt"
)

type Body_0x1203 struct {
	Appid    string
	Filename uint32
	Filelen  uint32
}

func (entity *Body_0x1203) MsgID() MsgID {
	return Msg_0x1203
}

func (entity *Body_0x1203) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Appid, 20)

	writer.WriteUint32(entity.Filename)

	writer.WriteUint32(entity.Filelen)

	return writer.Bytes(), nil
}

func (entity *Body_0x1203) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Appid, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.Filename, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	entity.Filelen, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
