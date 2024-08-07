package th

import (
	_ "bytes"
	_ "encoding/hex"
	"fmt"
	"strings"
)

type Message struct {
    Body Entity
}

func (message *Message) Encode() ([]byte, error) {
    body, err := message.Body.Encode()
    if err != nil {
        return nil, err
    }

    return body, nil
}

// 协议解码
func (message *Message) Decode(data []byte) error {
    entity, _, err := message.decodeBody(data[:])
    if err == nil {
        message.Body = entity
    } else {
        fmt.Println("decode err", err)
    }

    return nil
}
func (message *Message) decodeBody(data []byte) (Entity, int, error) {
    strData := string(data)
    strList := strings.Split(strData, "|")
    strTyp := strList[0]

    creator, ok := entityMapper[strTyp]
    if !ok {
        return nil, 0, ErrTypeNotRegistered
    }
    entity := creator()

    count, err := entity.Decode(data)
    if err != nil {
        return nil, 0, err
    }
    return entity, count, nil
}
