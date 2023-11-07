package kks

import (
    "encoding/binary"
    "encoding/json"
    "fmt"
    "golang.org/x/text/encoding/unicode"
    "golang.org/x/text/transform"
    "github.com/rayjay214/parser/common"
    _ "strconv"
)

type Kks_0x15 struct {
    Proto     uint8
    CmdLen    uint8
    SysFlag   uint32
    Encodings uint8
    Content   string
}

func (entity *Kks_0x15) MsgID() MsgID {
    return Msg_0x15
}

func (entity *Kks_0x15) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *Kks_0x15) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.Proto, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.CmdLen, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.SysFlag, err = reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    var encoding uint16
    offset := 2 + entity.CmdLen
    encoding = binary.BigEndian.Uint16(data[offset : offset+2])

    if encoding == 1 {
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

    entity.Encodings, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}

func (entity Kks_0x15) MarshalJSON() ([]byte, error) {
    type Alias Kks_0x15

    type New0x15 struct {
        Proto string
        Alias
    }

    s := New0x15{
        Alias: Alias(entity),
    }

    s.Proto = "0x" + fmt.Sprintf("%02x", entity.Proto)

    return json.Marshal(s)
}
