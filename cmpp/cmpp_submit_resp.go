package cmpp

import (
    _ "encoding/json"
    _ "fmt"
    "parser/common"
    _ "strconv"
)

type CMPP_SUBMIT_RESP struct {
    MsgSymbol uint64
    Result    uint32
}

func (entity *CMPP_SUBMIT_RESP) MsgID() MsgID {
    return MSG_CMPP_SUBMIT_RESP
}

func (entity *CMPP_SUBMIT_RESP) Encode() ([]byte, error) {
    //todo
    return nil, nil
}

func (entity *CMPP_SUBMIT_RESP) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.MsgSymbol, err = reader.ReadUint64()
    if err != nil {
        return 0, err
    }

    entity.Result, err = reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
