package jt808

import (
    "parser/common"
    "golang.org/x/text/encoding/unicode"
    "golang.org/x/text/transform"
)

type T808_0x1300 struct {
    AckSeqNo uint16
    Content  string
}

func (entity *T808_0x1300) MsgID() MsgID {
    return MsgT808_0x1300
}

func (entity *T808_0x1300) Encode() ([]byte, error) {
    return nil, nil
}

func (entity *T808_0x1300) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error

    entity.AckSeqNo, err = reader.ReadUint16()
    if err != nil {
        return 0, err
    }

    uniCodeContent, err := reader.Read()
    if err != nil {
        return 0, err
    }
    utf8Content, _, _ := transform.Bytes(unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), uniCodeContent)
    entity.Content = string(utf8Content)

    return 0, nil
}
