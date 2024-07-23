package ipc

import (
    _ "fmt"
)

type StateInfo struct {
	Key   uint16
	Value uint16
}

type Body_0x0004 struct {
	Num       uint8
	StateList []StateInfo
}

func (entity *Body_0x0004) MsgID() MsgID {
	return Msg_0x0004
}

func (entity *Body_0x0004) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteByte(entity.Num)

	for i := 0; i < int(entity.Num); i++ {
		writer.WriteUint16(entity.StateList[i].Key)
		writer.WriteUint16(entity.StateList[i].Value)
	}

	return writer.Bytes(), nil
}

func (entity *Body_0x0004) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error

	entity.Num, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(entity.Num); i++ {
		var info StateInfo
		info.Key, err = reader.ReadUint16()
		if err != nil {
			return 0, err
		}
		info.Value, err = reader.ReadUint16()
		if err != nil {
			return 0, err
		}
		entity.StateList = append(entity.StateList, info)
	}

	return len(data) - reader.Len(), nil
}
