package ipc

// 消息头
type Header struct {
    Prefix uint16
    MsgId  MsgID
    MsgLen uint16
    Seq    uint16
    UidLen uint16
    Uid    string
}

// 协议编码
func (header *Header) Encode() ([]byte, error) {
    writer := NewWriter()

    writer.WriteUint16(uint16(0x8686))

    writer.WriteUint16(uint16(header.MsgId))

    writer.WriteUint16(header.MsgLen)

    writer.WriteUint16(header.Seq)

    writer.WriteUint16(header.UidLen)

    writer.WriteString(header.Uid, int(header.UidLen))

    return writer.Bytes(), nil
}

// 协议解码
func (header *Header) Decode(data []byte) error {
    if len(data) < MessageHeaderSize {
        return ErrInvalidHeader
    }
    reader := NewReader(data)

    var err error
    header.Prefix, err = reader.ReadUint16()
    if err != nil {
        return ErrInvalidHeader
    }

    var msgId uint16
    msgId, err = reader.ReadUint16()
    header.MsgId = MsgID(msgId)
    if err != nil {
        return ErrInvalidHeader
    }

    header.MsgLen, err = reader.ReadUint16()
    if err != nil {
        return ErrInvalidHeader
    }

    header.Seq, err = reader.ReadUint16()
    if err != nil {
        return ErrInvalidHeader
    }

    header.UidLen, err = reader.ReadUint16()
    if err != nil {
        return ErrInvalidHeader
    }

    header.Uid, err = reader.ReadString(int(header.UidLen))
    if err != nil {
        return ErrInvalidHeader
    }
    return nil
}
