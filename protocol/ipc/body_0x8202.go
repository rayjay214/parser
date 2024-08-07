package ipc

import (
	_ "fmt"
)

type Body_0x8202 struct {
	Appid    string
	Filename uint32
}

func (entity *Body_0x8202) MsgID() MsgID {
	return Msg_0x8202
}

func (entity *Body_0x8202) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Appid, 20)

	writer.WriteUint32(entity.Filename)

	return writer.Bytes(), nil
}

func (entity *Body_0x8202) Decode(data []byte) (int, error) {
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

	return len(data) - reader.Len(), nil
}
