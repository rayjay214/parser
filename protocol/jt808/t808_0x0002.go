package jt808

import (
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/rayjay214/parser/protocol/jt808/errors"
	"github.com/rayjay214/parser/protocol/jt808/extra"
	log "github.com/sirupsen/logrus"
)

// 终端心跳
type T808_0x0002 struct {
	Extras []extra.Entity
}

func (entity *T808_0x0002) MsgID() MsgID {
	return MsgT808_0x0002
}

func (entity *T808_0x0002) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T808_0x0002) Decode(data []byte) (int, error) {
	//解析扩展协议
	reader := common.NewReader(data)
	if len(data) > 2 {
		extras := make([]extra.Entity, 0)
		buffer := data[len(data)-reader.Len():]
		for {
			if len(buffer) < 2 {
				break
			}
			id, length := buffer[0], int(buffer[1])
			buffer = buffer[2:]
			if len(buffer) < length {
				return 0, errors.ErrInvalidExtraLength
			}

			extraEntity, count, err := extra.Decode(id, buffer[:length])
			if err != nil {
				if err == errors.ErrTypeNotRegistered {
					buffer = buffer[length:]
					log.WithFields(log.Fields{
						"id": fmt.Sprintf("0x%x", id),
					}).Warn("[JT/T808] unknown T808_0x0200 extra type")
					continue
				}
				return 0, err
			}
			if count != length {
				return 0, errors.ErrInvalidExtraLength
			}
			extras = append(extras, extraEntity)
			buffer = buffer[length:]
		}
		if len(extras) > 0 {
			entity.Extras = extras
		}
		return len(data) - reader.Len(), nil
	}

	return 0, nil
}

func (entity T808_0x0002) MarshalJSON() ([]byte, error) {
	type Alias T808_0x0002

	type New0002 struct {
		Alias
		Extras map[string]interface{}
	}

	s := New0002{
		Alias:  Alias(entity),
		Extras: map[string]interface{}{},
	}

	for _, v := range entity.Extras {
		var strId string
		var val interface{}
		strId = "0x" + fmt.Sprintf("%02x", v.ID())
		switch vPrint := v.ToPrint().(type) {
		case map[string]interface{}:
			val = vPrint
		default:
			fmt.Println("%T", vPrint)
		}
		s.Extras[strId] = val
	}
	return json.Marshal(s)
}
