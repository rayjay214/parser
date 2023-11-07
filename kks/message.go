package kks

import (
    "bytes"
    _ "encoding/hex"
    "fmt"
)

type Message_0x79 struct {
    Header Header_0x79
    Body   Entity
}

type Message_0x78 struct {
    Header Header_0x78
    Body   Entity
}

func (message *Message_0x79) Encode() ([]byte, error) {
    body, err := message.Body.Encode()
    if err != nil {
        return nil, err
    }

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
func (message *Message_0x79) Decode(data []byte) error {
    if len(data) < 2 || (data[0] != 0x79 || data[1] != 0x79) {
        return ErrInvalidMessage
    }

    //todo check checksum

    if len(data) < MessageHeader79Size {
        return ErrInvalidHeader
    }
    var header Header_0x79
    err := header.Decode(data)
    if err != nil {
        return err
    }

    entity, _, err := message.decodeBody(data[MessageHeader79Size:])
    if err == nil {
        message.Body = entity
    } else {
        fmt.Println("decode err", err)
    }

    message.Header = header
    return nil
}
func (message *Message_0x79) decodeBody(data []byte) (Entity, int, error) {
    typ := data[0]
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

func (message *Message_0x79) write(buffer *bytes.Buffer, data []byte) *Message_0x79 {
    for _, b := range data {
        buffer.WriteByte(b)
    }
    return message
}

func (message *Message_0x79) computeChecksum(data []byte, checkSum byte) byte {
    for _, b := range data {
        checkSum = checkSum ^ b
    }
    return checkSum
}

func (message *Message_0x78) Encode() ([]byte, error) {
    body, err := message.Body.Encode()
    if err != nil {
        return nil, err
    }

    message.Header.MsgLen = byte(len(body))

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

func (message *Message_0x78) Decode(data []byte) error {
    if len(data) < 2 || (data[0] != 0x78 || data[1] != 0x78) {
        return ErrInvalidMessage
    }

    //todo check checksum

    if len(data) < MessageHeader78Size {
        return ErrInvalidHeader
    }
    var header Header_0x78
    err := header.Decode(data)
    if err != nil {
        return err
    }

    message.Header = header

    entity, _, err := message.decodeBody(data[MessageHeader78Size:])
    if err == nil {
        message.Body = entity
    } else {
        fmt.Println("decode err", err)
    }

    return nil
}
func (message *Message_0x78) decodeBody(data []byte) (Entity, int, error) {
    typ := data[0]
    //几个不同的消息用了同一个proto, 这里要特殊处理
    if typ == 0x17 && message.Header.MsgLen != 36 {
        fmt.Printf("%q", data[7])
        if data[7] == 'D' { //ADDRESS
            typ = 0xF0
        } else if data[7] == 'L' { //ALARMSMS
            typ = 0xF1
        } else {
            return nil, 0, ErrTypeNotRegistered
        }
    }
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

func (message *Message_0x78) write(buffer *bytes.Buffer, data []byte) *Message_0x78 {
    for _, b := range data {
        buffer.WriteByte(b)
    }
    return message
}

func (message *Message_0x78) computeChecksum(data []byte, checkSum byte) byte {
    for _, b := range data {
        checkSum = checkSum ^ b
    }
    return checkSum
}
