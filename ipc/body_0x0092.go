package ipc

import (
	"bytes"
	_ "fmt"
	"io"
	"io/ioutil"
)

type Body_0x0092 struct {
	SeqNo  uint16
	Packet io.Reader //数据包
}

func (entity *Body_0x0092) MsgID() MsgID {
	return Msg_0x0092
}

func (entity *Body_0x0092) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteUint16(entity.SeqNo)

	if entity.Packet != nil {
		data, err := ioutil.ReadAll(entity.Packet)
		if err != nil {
			return nil, err
		}
		writer.Write(data)
	}

	return writer.Bytes(), nil
}

func (entity *Body_0x0092) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error

	entity.SeqNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity *Body_0x0092) DecodePacket(data []byte) error {
	n, err := entity.Decode(data)
	if err != nil {
		return err
	}

	//checksum和结束符不能读

	entity.Packet = bytes.NewReader(data[n : len(data)-2])
	return nil
}

func (entity *Body_0x0092) GetTag() uint32 {
	return 0
}

func (entity *Body_0x0092) GetReader() io.Reader {
	return entity.Packet
}

func (entity *Body_0x0092) SetReader(reader io.Reader) {
	entity.Packet = reader
}
