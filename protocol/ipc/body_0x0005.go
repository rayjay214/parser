package ipc

import (
	_ "fmt"
)

type StateInfoExt struct {
	Key   uint16
	Len   uint16
	Value interface{}
}

type Body_0x0005 struct {
	Num       uint8
	StateList []StateInfoExt
}

func (entity *Body_0x0005) MsgID() MsgID {
	return Msg_0x0005
}

func (entity *Body_0x0005) Encode() ([]byte, error) {
	writer := NewWriter()

	return writer.Bytes(), nil
}

func (entity *Body_0x0005) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error

	entity.Num, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(entity.Num); i++ {
		var info StateInfoExt
		info.Key, err = reader.ReadUint16()
		if err != nil {
			return 0, err
		}
		info.Len, err = reader.ReadUint16()
		if err != nil {
			return 0, err
		}
		switch info.Key {
		case 0x01:
			info.Value, err = reader.ReadUint32()
		case 0x02:
			info.Value, err = reader.ReadUint32()
		}
		if err != nil {
			return 0, err
		}

		entity.StateList = append(entity.StateList, info)
	}

	return len(data) - reader.Len(), nil
}
