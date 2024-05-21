package ipc

import (
	_ "fmt"
)

type SdFile struct {
	BeginTime uint32
	Cos       uint16
}

type Body_0x1201 struct {
	Appid    string
	Seq      uint16
	FileType string
	FileNum  uint16
	Files    []SdFile
}

func (entity *Body_0x1201) MsgID() MsgID {
	return Msg_0x1201
}

func (entity *Body_0x1201) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Appid, 20)

	writer.WriteUint16(entity.Seq)

	writer.WriteString(entity.FileType, 8)

	writer.WriteUint16(entity.FileNum)

	for i := 0; i < int(entity.FileNum); i++ {
		writer.WriteUint32(entity.Files[i].BeginTime)
		writer.WriteUint16(entity.Files[i].Cos)
	}

	return writer.Bytes(), nil
}

func (entity *Body_0x1201) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Appid, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.Seq, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	entity.FileType, err = reader.ReadString(8)
	if err != nil {
		return 0, err
	}

	entity.FileNum, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(entity.FileNum); i++ {
		var f SdFile
		f.BeginTime, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		f.Cos, err = reader.ReadUint16()
		if err != nil {
			return 0, err
		}
		entity.Files = append(entity.Files, f)
	}

	return len(data) - reader.Len(), nil
}
