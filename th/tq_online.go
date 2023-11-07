package th

import (
    _ "encoding/json"
    _ "fmt"
    "github.com/rayjay214/parser/common"
    _ "strconv"
    _ "strings"
)

type TH_ONLINE struct {
    Proto string
}

func (entity *TH_ONLINE) MsgID() MsgID {
    return Msg_ONLINE
}

func (entity *TH_ONLINE) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *TH_ONLINE) Decode(data []byte) (int, error) {

    //todo

    return len(data), nil
}
