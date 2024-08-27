package jt808_base

import (
	"bytes"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"github.com/rayjay214/link"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/rayjay214/parser/protocol/jt808"
	"github.com/rayjay214/parser/protocol/jt808/errors"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"io"
)

type Protocol struct {
	privateKey *rsa.PrivateKey
}

// 创建编解码器
func (p Protocol) NewCodec(rw io.ReadWriter) (link.Codec, error) {
	codec := &ProtocolCodec{
		w:               rw,
		r:               rw,
		privateKey:      p.privateKey,
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
	publicKey       *rsa.PublicKey
	privateKey      *rsa.PrivateKey
	bufferReceiving *bytes.Buffer
}

// 获取RSA公钥
func (codec *ProtocolCodec) GetPublicKey() *rsa.PublicKey {
	return codec.publicKey
}

// 设置RSA公钥
func (codec *ProtocolCodec) SetPublicKey(publicKey *rsa.PublicKey) {
	codec.publicKey = publicKey
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
	message, ok := msg.(jt808.Message)
	if !ok {
		log.WithFields(log.Fields{
			"reason": errors.ErrInvalidMessage,
		}).Error("[JT/T 808] failed to write message")
		return errors.ErrInvalidMessage
	}

	var err error
	var data []byte
	if codec.publicKey == nil || !message.Header.Property.IsEnableEncrypt() {
		data, err = message.Encode()
	} else {
		data, err = message.Encode(codec.publicKey)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("0x%x", message.Header.MsgID),
			"reason": err,
		}).Error("[JT/T 808] failed to write message")
		return err
	}

	_, err = codec.w.Write(data)
	if err != nil {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("0x%x", message.Header.MsgID),
			"reason": err,
		}).Error("[JT/T 808] failed to write message")
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
func (codec *ProtocolCodec) readFromBuffer() (jt808.Message, bool, error) {
	if codec.bufferReceiving.Len() == 0 {
		return jt808.Message{}, false, nil
	}

	data := codec.bufferReceiving.Bytes()
	if data[0] != jt808.PrefixID {
		i := 0
		for ; i < len(data); i++ {
			if data[i] == jt808.PrefixID {
				break
			}
		}
		codec.bufferReceiving.Next(i)
		log.WithFields(log.Fields{
			"data":   hex.EncodeToString(data),
			"reason": errors.ErrNotFoundPrefixID,
		}).Error("[JT/T 808] failed to receive message")
		return jt808.Message{}, false, errors.ErrNotFoundPrefixID
	}

	end := 1
	for ; end < len(data); end++ {
		if data[end] == jt808.PrefixID {
			break
		}
	}
	if end == len(data) {
		return jt808.Message{}, false, nil
	}

	var message jt808.Message
	if err := message.Decode(data[:end+1], codec.privateKey); err != nil {
		codec.bufferReceiving.Next(end + 1)
		log.WithFields(log.Fields{
			"data":   fmt.Sprintf("%x", hex.EncodeToString(data[:end+1])),
			"reason": err,
		}).Error("[JT/T 808] failed to receive message")
		return jt808.Message{}, false, err
	}
	codec.bufferReceiving.Next(end + 1)

	if message.Header.MsgID != jt808.MsgT808_0x0118 {
		storage.RawLogChannel <- storage.LogRow{
			Imei:      message.Header.Imei,
			Direction: "R",
			Message:   fmt.Sprintf("%x", common.GetHex(data[:end+1])),
		}
	}

	return message, true, nil
}
