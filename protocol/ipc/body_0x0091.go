package ipc

import (
	_ "fmt"
)

type Body_0x0091 struct {
	LatestVersion string
	FileLen       uint32
	Crc           uint32
	OtaIp         string
	OtaPort       uint16
}

func (entity *Body_0x0091) MsgID() MsgID {
	return Msg_0x0091
}

func (entity *Body_0x0091) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.LatestVersion, 32)

	writer.WriteUint32(entity.FileLen)

	writer.WriteUint32(entity.Crc)

	writer.WriteString(entity.OtaIp, 32)

	writer.WriteUint16(entity.OtaPort)

	return writer.Bytes(), nil
}

func (entity *Body_0x0091) Decode(data []byte) (int, error) {
	reader := NewReader(data)
	var err error

	entity.LatestVersion, err = reader.ReadString(32)
	if err != nil {
		return 0, err
	}

	entity.FileLen, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	entity.Crc, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	entity.OtaIp, err = reader.ReadString(32)
	if err != nil {
		return 0, err
	}

	entity.OtaPort, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
