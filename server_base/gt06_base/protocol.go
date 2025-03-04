package gt06_base

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/rayjay214/link"
	"github.com/rayjay214/parser/protocol/gt06"
	"github.com/rayjay214/parser/protocol/jt808/errors"
	log "github.com/sirupsen/logrus"
	"io"
)

type Protocol struct {
}

// 创建编解码器
func (p Protocol) NewCodec(rw io.ReadWriter) (link.Codec, error) {
	codec := &ProtocolCodec{
		w:               rw,
		r:               rw,
		bufferReceiving: bytes.NewBuffer(nil),
	}
	codec.closer, _ = rw.(io.Closer)
	return codec, nil
}

// 编解码器
type ProtocolCodec struct {
	w               io.Writer
	r               io.Reader
	closer          io.Closer
	bufferReceiving *bytes.Buffer
}

// 关闭读写
func (codec *ProtocolCodec) Close() error {
	if codec.closer != nil {
		return codec.closer.Close()
	}
	return nil
}

// 发送消息
func (codec *ProtocolCodec) Send(msg interface{}) error {
	message, ok := msg.(gt06.Message)
	if !ok {
		log.WithFields(log.Fields{
			"reason": errors.ErrInvalidMessage,
		}).Error("[gt06] failed to write message")
		return errors.ErrInvalidMessage
	}

	var err error
	var data []byte
	data, err = message.Encode()

	if err != nil {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("%v", message.Body.MsgID),
			"reason": err,
		}).Error("[gt06] failed to write message")
		return err
	}

	_, err = codec.w.Write(data)
	if err != nil {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("0x%v", message.Body.MsgID),
			"reason": err,
		}).Error("[gt06] failed to write message")
		return err
	}

	return nil
}

// 接收消息
func (codec *ProtocolCodec) Receive() (interface{}, error) {
	message, ok, err := codec.readFromBuffer()
	if ok {
		return message, nil
	}
	if err != nil {
		log.Warnf("receive err %v", err)
		return nil, err
	}

	var buffer [128]byte
	for {
		count, err := io.ReadAtLeast(codec.r, buffer[:], 1)
		if err != nil {
			log.Warnf("receive err %v", err)
			return nil, err
		}
		codec.bufferReceiving.Write(buffer[:count])

		if codec.bufferReceiving.Len() == 0 {
			continue
		}
		if codec.bufferReceiving.Len() > 0xffff {
			return nil, errors.ErrBodyTooLong
		}

		message, ok, err := codec.readFromBuffer()
		if ok {
			return message, nil
		}
		if err != nil {
			log.Warnf("receive err %v", err)
			return nil, err
		}
	}
}

// 从缓冲区读取
func (codec *ProtocolCodec) readFromBuffer() (gt06.Message, bool, error) {
	if codec.bufferReceiving.Len() == 0 {
		return gt06.Message{}, false, nil
	}

	data := codec.bufferReceiving.Bytes()

	if len(data) < 4 {
		return gt06.Message{}, false, nil
	}

	prefix := binary.BigEndian.Uint16(data[:2])
	if prefix != 0x7878 && prefix != 0x7979 {
		return gt06.Message{}, false, errors.ErrInvalidHeader
	}

	var dataLen, msgLen int
	if prefix == 0x7979 {
		dataLen = int(binary.BigEndian.Uint16(data[2:4]))
		msgLen = dataLen + 6 //起始位+长度+停止位
	} else {
		dataLen = int(data[2])
		msgLen = dataLen + 5 //起始位+长度+停止位
	}

	if len(data) < msgLen {
		return gt06.Message{}, false, nil
	}

	var message gt06.Message
	if err := message.Decode(data[:msgLen]); err != nil {
		codec.bufferReceiving.Next(msgLen)
		log.WithFields(log.Fields{
			"data":   fmt.Sprintf("0x%x", hex.EncodeToString(data[:msgLen])),
			"reason": err,
		}).Error("failed to receive message")
		return gt06.Message{}, false, err
	}
	codec.bufferReceiving.Next(msgLen)

	return message, true, nil
}
