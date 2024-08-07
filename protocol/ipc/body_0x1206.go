package ipc

import (
	_ "fmt"
)

type Result struct {
	Date    uint32
	HasFile uint8
}

type Body_0x1206 struct {
	Appid   string
	Num     uint8
	Summary []Result
}

func (entity *Body_0x1206) MsgID() MsgID {
	return Msg_0x1206
}

func (entity *Body_0x1206) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Appid, 20)

	writer.WriteByte(entity.Num)

	for i := 0; i < int(entity.Num); i++ {
		writer.WriteUint32(entity.Summary[i].Date)
		writer.WriteByte(entity.Summary[i].HasFile)
	}

	return writer.Bytes(), nil
}

func (entity *Body_0x1206) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Appid, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.Num, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(entity.Num); i++ {
		var f Result
		f.Date, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		f.HasFile, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}
		entity.Summary = append(entity.Summary, f)
	}

	return len(data) - reader.Len(), nil
}
