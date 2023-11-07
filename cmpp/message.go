package cmpp

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
    //todo
    return nil, nil
}

// 协议解码
func (message *Message) Decode(data []byte) error {
    //todo check

    var header Header
    err := header.Decode(data)
    if err != nil {
        return err
    }

    entity, _, err := message.decodeBody(uint32(header.CmdId), data[12:])
    if err == nil {
        message.Body = entity
    } else {
        fmt.Println("decode err", err)
    }

    message.Header = header
    return nil
}
func (message *Message) decodeBody(typ uint32, data []byte) (Entity, int, error) {
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
