package gt06

import (
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/common"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	_ "strconv"
)

type Kks_0x21 struct {
	Proto     uint8
	SysFlag   uint32
	Encodings uint8
	Content   string
	SeqNo     uint16
}

func (entity *Kks_0x21) GetSeqNo() uint16 {
	return entity.SeqNo
}

func (entity *Kks_0x21) MsgID() MsgID {
	return Msg_0x21
}

func (entity *Kks_0x21) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *Kks_0x21) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error
	entity.Proto, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.SysFlag, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	entity.Encodings, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	if entity.Encodings == 2 {
		var uniCodeContent []byte
		uniCodeContent, err = reader.Read(len(data) - 12)
		if err != nil {
			return 0, err
		}

		utf8Content, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), uniCodeContent)
		entity.Content = string(utf8Content)
	} else {
		entity.Content, err = reader.ReadString(len(data) - 12)
		if err != nil {
			return 0, err
		}
	}

	entity.SeqNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity Kks_0x21) MarshalJSON() ([]byte, error) {
	type Alias Kks_0x21

	type New0x21 struct {
		Proto string
		Alias
	}

	s := New0x21{
		Alias: Alias(entity),
	}

	s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

	return json.Marshal(s)
}
