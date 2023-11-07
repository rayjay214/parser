package th

import (
    _ "encoding/json"
    _ "fmt"
    "github.com/rayjay214/parser/common"
    _ "strconv"
    _ "strings"
)

type TH_FOTA struct {
    Proto string
}

func (entity *TH_FOTA) MsgID() MsgID {
    return Msg_FOTA
}

func (entity *TH_FOTA) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *TH_FOTA) Decode(data []byte) (int, error) {

    //todo

    return len(data), nil
}
