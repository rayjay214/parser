package cmpp

import (
    _ "encoding/json"
    _ "fmt"
    "parser/common"
    _ "strconv"
)

type CMPP_CONNECT_RESP struct {
    Status            uint32
    AuthenticatorISMG string
    Version           byte
}

func (entity *CMPP_CONNECT_RESP) MsgID() MsgID {
    return MSG_CMPP_CONNECT_RESP
}

func (entity *CMPP_CONNECT_RESP) Encode() ([]byte, error) {
    //todo
    return nil, nil
}

func (entity *CMPP_CONNECT_RESP) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.Status, err = reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    entity.AuthenticatorISMG, err = reader.ReadString(16)
    if err != nil {
        return 0, err
    }

    entity.Version, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
