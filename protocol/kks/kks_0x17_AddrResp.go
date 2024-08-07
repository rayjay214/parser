package kks

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"github.com/rayjay214/parser/protocol/common"
	_ "strconv"
)

type Kks_0x17_AddrResp struct {
	Proto       uint8
	ContentLen  uint8
	AlarmFlag   uint32
	AddrContent string
	Phone       string
}

func (entity *Kks_0x17_AddrResp) MsgID() MsgID {
	return Msg_0x17_AddrResp
}

func (entity *Kks_0x17_AddrResp) Encode() ([]byte, error) {
	writer := common.NewWriter()

	//todo

	return writer.Bytes(), nil
}

func (entity *Kks_0x17_AddrResp) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error
	entity.Proto, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.ContentLen, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.AlarmFlag, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	_, err = reader.ReadString(7) //ADDRESS
	if err != nil {
		return 0, err
	}

	_, err = reader.ReadString(2) //&&
	if err != nil {
		return 0, err
	}

	addrLen := int(entity.ContentLen) - 4 - 7 - 2 - 2 - 21 - 2
	var uniCodeContent []byte
	uniCodeContent, err = reader.Read(int(addrLen))
	if err != nil {
		return 0, err
	}

	utf8Content, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), uniCodeContent)
	entity.AddrContent = string(utf8Content)

	_, err = reader.ReadString(2) //&&
	if err != nil {
		return 0, err
	}

	entity.Phone, err = reader.ReadString(21)
	if err != nil {
		return 0, err
	}

	_, err = reader.ReadString(2) //##
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity Kks_0x17_AddrResp) MarshalJSON() ([]byte, error) {
	type Alias Kks_0x17_AddrResp

	type New0x17_AddrResp struct {
		Proto     string
		AlarmFlag string
		Alias
	}

	s := New0x17_AddrResp{
		Alias: Alias(entity),
	}

	s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)
	s.AlarmFlag = "0x" + fmt.Sprintf("%02x", entity.AlarmFlag)

	return json.Marshal(s)
}
