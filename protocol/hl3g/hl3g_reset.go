package hl3g

import (
	_ "encoding/json"
	_ "fmt"
	"github.com/rayjay214/parser/protocol/common"
	_ "github.com/shopspring/decimal"
	_ "strconv"
	_ "time"
)

type HL3G_RESET struct {
}

func (entity *HL3G_RESET) MsgID() MsgID {
	return Msg_RESET
}

func (entity *HL3G_RESET) Encode() ([]byte, error) {
	writer := common.NewWriter()

	return writer.Bytes(), nil
}

func (entity *HL3G_RESET) Decode(data []byte) (int, error) {
	return len(data), nil
}
