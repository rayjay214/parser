package ipc

import (
	_ "fmt"
)

type Body_0x0012 struct {
	PkgSize uint16
	PkgSeq  uint16
	Version string
}

func (entity *Body_0x0012) MsgID() MsgID {
	return Msg_0x0012
}

func (entity *Body_0x0012) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteUint16(entity.PkgSize)

	writer.WriteUint16(entity.PkgSeq)

	writer.WriteString(entity.Version, 32)

	return writer.Bytes(), nil
}

func (entity *Body_0x0012) Decode(data []byte) (int, error) {
	reader := NewReader(data)
	var err error

	entity.PkgSize, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.PkgSeq, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.Version, err = reader.ReadString(32)
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
