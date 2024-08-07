package ipc

import (
	_ "fmt"
)

type Body_0x8207 struct {
	Appid    string
	FileName uint32
}

func (entity *Body_0x8207) MsgID() MsgID {
	return Msg_0x8207
}

func (entity *Body_0x8207) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Appid, 20)

	writer.WriteUint32(entity.FileName)

	return writer.Bytes(), nil
}

func (entity *Body_0x8207) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Appid, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.FileName, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
