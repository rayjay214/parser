package hl3g

import (
	_ "bytes"
	_ "encoding/hex"
	"fmt"
)

type Message struct {
	Header Header
	Body   Entity
}

func (message *Message) Encode() ([]byte, error) {
	body, err := message.Body.Encode()
	if err != nil {
		return nil, err
	}

	return body, nil
}

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
	var strTyp string
	if data[0] == '*' { //ASCII
		buf := data[15:17]
		strTyp = string(buf)
	} else { //BINARY
		strTyp = "Normal"
	}

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
