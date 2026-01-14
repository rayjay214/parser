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

// todo
func (message *Message) Encode() ([]byte, error) {
	var body []byte
	var err error
	if message.Body != nil {
		body, err = message.Body.Encode()
		if err != nil {
			return nil, err
		}
	}

	header, err := message.Header.Encode()
	if err != nil {
		return nil, err
	}

	if message.Body == nil {
		header = header[:len(header)-1]
	}

	msg := string(header) + string(body) + "]"

	return []byte(msg), nil
}

func (message *Message) Decode(data []byte) error {
	strData := string(data)
	if strData[0:3] != "[3G" {
		return ErrInvalidMessage
	}

	var h Header
	err := h.Decode(data)
	if err != nil {
		return err
	}
	headerLen := len(h.Prefix) + len(h.Imei) + len(h.MsgLen) + len(h.Proto) + 3
	message.Header = h

	entity, _, err := message.decodeBody(data[headerLen:])
	if err == nil {
		message.Body = entity
	} else {
		fmt.Println("decode err", err)
	}

	return nil
}

func (message *Message) decodeBody(data []byte) (Entity, int, error) {
	creator, ok := entityMapper[message.Header.Proto]
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
