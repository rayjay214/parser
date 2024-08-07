package ipc

import (
	_ "fmt"
)

type Body_0x1301 struct {
	Seq     uint16
	Len     uint8
	Content string
}

func (entity *Body_0x1301) MsgID() MsgID {
	return Msg_0x1301
}

func (entity *Body_0x1301) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteUint16(entity.Seq)

	writer.WriteByte(entity.Len)

	writer.WriteString(entity.Content, int(entity.Len))

	return writer.Bytes(), nil
}

func (entity *Body_0x1301) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Seq, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Len, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.Content, err = reader.ReadString(int(entity.Len))
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
