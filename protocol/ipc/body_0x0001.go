package ipc

import (
	_ "fmt"
)

type Body_0x0001 struct {
	Version string
}

func (entity *Body_0x0001) MsgID() MsgID {
	return Msg_0x0001
}

func (entity *Body_0x0001) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Version, 32)

	return writer.Bytes(), nil
}

func (entity *Body_0x0001) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	if len(data) > 30 {
		var err error
		entity.Version, err = reader.ReadString(32)
		if err != nil {
			return 0, err
		}
	}

	return len(data) - reader.Len(), nil
}
