package ota

import (
    "bytes"
    _ "encoding/hex"
    "fmt"
)

type Message struct {
    Header Header
    Body   Entity
}

func (message *Message) EncodeEmpty91() ([]byte, error) {
    var body = []byte{0x00}

    message.Header.MsgId = MsgID(0x91)
    message.Header.MsgLen = uint16(len(body))

    header, err := message.Header.Encode()
    if err != nil {
        return nil, err
    }
    checkSum := byte(0x00)
    checkSum = message.computeChecksum(header, checkSum)
    checkSum = message.computeChecksum(body, checkSum)
    tail := byte(0x0d)

    buffer := bytes.NewBuffer(nil)
    buffer.Grow(len(header) + len(body) + 2)
    message.write(buffer, header).write(buffer, body).write(buffer, []byte{checkSum}).write(buffer, []byte{tail})
    return buffer.Bytes(), nil
}

func (message *Message) Encode() ([]byte, error) {
    body, err := message.Body.Encode()
    if err != nil {
        return nil, err
    }

    message.Header.MsgId = message.Body.MsgID()
    message.Header.MsgLen = uint16(len(body))

    header, err := message.Header.Encode()
    if err != nil {
        return nil, err
    }
    checkSum := byte(0x00)
    checkSum = message.computeChecksum(header, checkSum)
    checkSum = message.computeChecksum(body, checkSum)
    tail := byte(0x0d)

    buffer := bytes.NewBuffer(nil)
    buffer.Grow(len(header) + len(body) + 2)
    message.write(buffer, header).write(buffer, body).write(buffer, []byte{checkSum}).write(buffer, []byte{tail})
    return buffer.Bytes(), nil
}

// 协议解码
func (message *Message) Decode(data []byte) error {
    if len(data) < 2 || (data[0] != PrefixID || data[1] != PrefixID) {
        return ErrInvalidMessage
    }

    //todo check checksum

    if len(data) < MessageHeaderSize {
        return ErrInvalidHeader
    }
    var header Header
    err := header.Decode(data)
    if err != nil {
        return err
    }

    entity, _, err := message.decodeBody(uint8(header.MsgId), data[MessageHeaderSize:])
    if err == nil {
        message.Body = entity
    } else {
        fmt.Println("decode err", err)
    }

    message.Header = header
    return nil
}
func (message *Message) decodeBody(typ uint8, data []byte) (Entity, int, error) {
    creator, ok := entityMapper[typ]
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

func (message *Message) write(buffer *bytes.Buffer, data []byte) *Message {
    for _, b := range data {
        buffer.WriteByte(b)
    }
    return message
}

func (message *Message) computeChecksum(data []byte, checkSum byte) byte {
    for _, b := range data {
        checkSum = checkSum ^ b
    }
    return checkSum
}
