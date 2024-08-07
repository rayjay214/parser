package cmpp

import (
	_ "encoding/json"
	_ "fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "strconv"
)

type CMPP_DELIVER_RESP struct {
	MsgSymbol uint64
	Result    uint8
}

func (entity *CMPP_DELIVER_RESP) MsgID() MsgID {
	return MSG_CMPP_DELIVER_RESP
}

func (entity *CMPP_DELIVER_RESP) Encode() ([]byte, error) {
	//todo
	return nil, nil
}

func (entity *CMPP_DELIVER_RESP) Decode(data []byte) (int, error) {
	reader := common.NewReader(data)

	var err error
	entity.MsgSymbol, err = reader.ReadUint64()
	if err != nil {
		return 0, err
	}

	entity.Result, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
