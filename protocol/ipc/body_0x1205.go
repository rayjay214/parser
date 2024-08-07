package ipc

import (
	_ "fmt"
)

type Body_0x1205 struct {
	Appid    string
	Result   uint8
	Filename uint32
}

func (entity *Body_0x1205) MsgID() MsgID {
	return Msg_0x1205
}

func (entity *Body_0x1205) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Appid, 20)

	writer.WriteByte(entity.Result)

	writer.WriteUint32(entity.Filename)

	return writer.Bytes(), nil
}

func (entity *Body_0x1205) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Appid, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.Result, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Filename, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
