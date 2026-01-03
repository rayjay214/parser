package hl3g

import (
	"github.com/rayjay214/parser/protocol/common"
	_ "strconv"
	"strings"
)

// 消息头
type Header struct {
	Prefix string
	Imei   string
	MsgLen string
	Proto  string
}

func (header *Header) Encode() ([]byte, error) {
	writer := common.NewWriter()

	writer.WriteString(header.Prefix)
	writer.WriteString("*")
	writer.WriteString(header.Imei)
	writer.WriteString("*")
	writer.WriteString(header.MsgLen)
	writer.WriteString("*")
	writer.WriteString(header.Proto)
	writer.WriteString(",")

	return writer.Bytes(), nil
}

func (header *Header) Decode(data []byte) error {
	strData := string(data)
	strList := strings.Split(strData, "*")
	if len(strList) != 4 {
		return ErrInvalidHeader
	}
	header.Prefix = strList[0]
	header.Imei = strList[1]
	header.MsgLen = strList[2]
	content := strList[3]
	if strings.Contains(content, ",") {
		strList2 := strings.Split(content, ",")
		header.Proto = strList2[0]
	} else {
		header.Proto = content
	}

	return nil
}
