package jt808

import (
	"crypto/rsa"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/rayjay214/parser/protocol/jt808/errors"
	"math/big"
)

// 终端 RSA公钥
type T808_0x0A00 struct {
	PublicKey *rsa.PublicKey
}

func (entity *T808_0x0A00) MsgID() MsgID {
	return MsgT808_0x0A00
}

func (entity *T808_0x0A00) Encode() ([]byte, error) {
	if entity.PublicKey == nil || entity.PublicKey.N == nil {
		return nil, errors.ErrInvalidPublicKey
	}
	if entity.PublicKey.Size() != 128 {
		return nil, errors.ErrInvalidPublicKey
	}

	writer := common.NewWriter()
	writer.WriteUint32(uint32(entity.PublicKey.E))
	writer.Write(entity.PublicKey.N.Bytes(), 128)
	return writer.Bytes(), nil
}

func (entity *T808_0x0A00) Decode(data []byte) (int, error) {
	if len(data) < 132 {
		return 0, errors.ErrInvalidBody
	}

	reader := common.NewReader(data)
	e, err := reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	n, err := reader.Read(128)
	if err != nil {
		return 0, err
	}

	entity.PublicKey = &rsa.PublicKey{
		E: int(e),
		N: big.NewInt(0).SetBytes(n),
	}
	return len(data) - reader.Len(), nil
}
