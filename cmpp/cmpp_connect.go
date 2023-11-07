package cmpp

import (
	_ "encoding/json"
	_ "fmt"
	"parser/common"
	_ "strconv"
)

type CMPP_CONNECT struct {
    SourceAddr          string
    AuthenticatorSource string
    Version             byte
    Timestamp           uint32
}

func (entity *CMPP_CONNECT) MsgID() MsgID {
    return MSG_CMPP_CONNECT
}

func (entity *CMPP_CONNECT) Encode() ([]byte, error) {
    //todo
    return nil, nil
}

func (entity *CMPP_CONNECT) Decode(data []byte) (int, error) {
    reader := common.NewReader(data)

    var err error
    entity.SourceAddr, err = reader.ReadString(6)
    if err != nil {
        return 0, err
    }

    entity.AuthenticatorSource, err = reader.ReadString(16)
    if err != nil {
        return 0, err
    }

    entity.Version, err = reader.ReadByte()
    if err != nil {
        return 0, err
    }

    entity.Timestamp, err = reader.ReadUint32()
    if err != nil {
        return 0, err
    }

    return len(data) - reader.Len(), nil
}
