package ipc

import (
	"bytes"
	_ "fmt"
	"io"
	"io/ioutil"
)

type Body_0x1204 struct {
	Appid  string
	SeqNo  uint16
	Packet io.Reader //数据包
}

func (entity *Body_0x1204) MsgID() MsgID {
	return Msg_0x1204
}

func (entity *Body_0x1204) Encode() ([]byte, error) {
	writer := NewWriter()

	writer.WriteString(entity.Appid, 20)

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

func (entity *Body_0x1204) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	entity.Appid, err = reader.ReadString(20)
	if err != nil {
		return 0, err
	}

	entity.SeqNo, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}

func (entity *Body_0x1204) DecodePacket(data []byte) error {
	n, err := entity.Decode(data)
	if err != nil {
		return err
	}

	//checksum和结束符不能读

	entity.Packet = bytes.NewReader(data[n : len(data)-2])
	return nil
}

func (entity *Body_0x1204) GetTag() uint32 {
	return 0
}

func (entity *Body_0x1204) GetReader() io.Reader {
	return entity.Packet
}

func (entity *Body_0x1204) SetReader(reader io.Reader) {
	entity.Packet = reader
}
