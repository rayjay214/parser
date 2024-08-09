package gmconf

import (
	_ "fmt"
	"github.com/rayjay214/parser/protocol/common"
)

type Body_0x13 struct {
	ConfType byte
}

type ConfInfo struct {
	ConfCat      byte //1:整数 2:小数 3:字符串 4:布尔型
	ConfName     string
	ConfLen      byte
	ContentInt   uint32
	ContentFloat uint32
	ContentStr   string
	ContentBool  byte
}

type Body_0x92 struct {
	HasConf  byte
	ConfList []ConfInfo
}

type Body_0x14 struct {
	ConfType byte
}

func (entity *Body_0x13) MsgID() MsgID {
	return Msg_0x13
}

func (entity *Body_0x13) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteByte(entity.ConfType)

	return writer.Bytes(), nil
}

func (entity *Body_0x13) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error
	entity.ConfType, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity *Body_0x14) MsgID() MsgID {
	return Msg_0x14
}

func (entity *Body_0x14) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteByte(entity.ConfType)

	return writer.Bytes(), nil
}

func (entity *Body_0x14) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error
	entity.ConfType, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity *Body_0x92) MsgID() MsgID {
	return Msg_0x92
}

func (entity *Body_0x92) Encode() ([]byte, error) {
	writer := common.NewWriter()

	return writer.Bytes(), nil
}

func (entity *Body_0x92) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error
	entity.HasConf, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	readLen := 1

	if entity.HasConf > 0 {
		for {
			left := len(data) - readLen
			if left <= 2 { //crc+0d
				break
			}
			var info ConfInfo
			info.ConfCat, err = reader.ReadByte()
			readLen += 1
			if err != nil {
				return 0, err
			}
			info.ConfName, err = reader.ReadString(10)
			readLen += 10
			switch info.ConfCat {
			case 1:
				info.ContentInt, err = reader.ReadUint32()
				readLen += 4
				if err != nil {
					return 0, err
				}
			case 2:
				info.ContentFloat, err = reader.ReadUint32()
				readLen += 4
				if err != nil {
					return 0, err
				}
			case 3:
				info.ConfLen, err = reader.ReadByte()
				readLen += 1
				if err != nil {
					return 0, err
				}
				info.ContentStr, err = reader.ReadString(int(info.ConfLen))
				if err != nil {
					return 0, err
				}
				readLen += int(info.ConfLen)
			case 4:
				info.ContentBool, err = reader.ReadByte()
				readLen += 1
				if err != nil {
					return 0, err
				}
			}
			entity.ConfList = append(entity.ConfList, info)
		}
	}

	return len(data) - reader.Len(), nil
}
