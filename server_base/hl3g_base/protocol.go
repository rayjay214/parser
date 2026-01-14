package hl3g_base

import (
	"bytes"
	"fmt"
	"github.com/rayjay214/link"
	"github.com/rayjay214/parser/protocol/hl3g"
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
	message, ok := msg.(hl3g.Message)
	if !ok {
		log.WithFields(log.Fields{
			"reason": errors.ErrInvalidMessage,
		}).Error("[hl3g] failed to write message")
		return errors.ErrInvalidMessage
	}

	var err error
	var data []byte
	data, err = message.Encode()

	if err != nil {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("%v", message.Header.Proto),
			"reason": err,
		}).Error("[hl3g] failed to write message")
		return err
	}

	_, err = codec.w.Write(data)
	if err != nil {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("%v", message.Header.Proto),
			"reason": err,
		}).Error("[hl3g] failed to write message")
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
func (codec *ProtocolCodec) readFromBuffer() (hl3g.Message, bool, error) {
	if codec.bufferReceiving.Len() == 0 {
		return hl3g.Message{}, false, nil
	}

	data := codec.bufferReceiving.Bytes()

	if len(data) < 3 {
		return hl3g.Message{}, false, nil
	}

	prefix := data[:3]
	if string(prefix) != "[3G" {
		return hl3g.Message{}, false, errors.ErrInvalidHeader
	}

	var msgLen int
	for _, b := range data {
		if b == ']' {
			msgLen += 1
			break
		}
		msgLen += 1
	}

	//缓冲区只有128字节，判断是否装下了整条消息
	if data[msgLen-1] != ']' {
		return hl3g.Message{}, false, nil
	}

	if len(data) < msgLen {
		return hl3g.Message{}, false, nil
	}

	var message hl3g.Message
	if err := message.Decode(data[:msgLen]); err != nil {
		codec.bufferReceiving.Next(msgLen)
		log.WithFields(log.Fields{
			"data":   fmt.Sprintf("%s", data[:msgLen]),
			"reason": err,
		}).Error("failed to receive message")
		return hl3g.Message{}, false, err
	}
	codec.bufferReceiving.Next(msgLen)

	log.Infof("hl3g recv msg %v", string(data[:msgLen]))

	return message, true, nil
}
