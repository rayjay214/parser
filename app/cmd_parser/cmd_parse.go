package main

import (
    "encoding/binary"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "github.com/rayjay214/parser/cmpp"
    "github.com/rayjay214/parser/jt808"
    "github.com/rayjay214/parser/kks"
    "github.com/rayjay214/parser/ota"
    "github.com/rayjay214/parser/th"
    "github.com/rayjay214/parser/tq"
    "os"
    "github.com/rayjay214/parser/common"
    "strings"
)

func main() {
    var data []byte
    var err interface{}
    if len(os.Args) < 2 {
        return
    }

    hexStr := os.Args[1]
    data, err = hex.DecodeString(hexStr)
    if err != nil {
        data = []byte(hexStr)
    }

    var out []byte
    prefix := data[0]
    switch prefix {
    case 0x7e:
        msg_id := binary.BigEndian.Uint16(data[1:])
        var transformed_data []byte
        if msg_id > 0x8000 {
            transformed_data = data
        } else {
            transformed_data = common.Transform808(data)
        }
        message := new(jt808.Message)
        err = message.Decode(transformed_data)
        out, _ = json.MarshalIndent(message, "", "   ")
    case 0x78:
        message := new(kks.Message_0x78)
        message.Decode(data)
        out, _ = json.MarshalIndent(message, "", "   ")
    case 0x79:
        message := new(kks.Message_0x79)
        message.Decode(data)
        out, _ = json.MarshalIndent(message, "", "   ")
    case 0x68:
        message := new(ota.Message)
        message.Decode(data)
        out, _ = json.MarshalIndent(message, "", "   ")
    case '*':
        message := new(tq.Message)
        message.Decode(data)
        out, _ = json.MarshalIndent(message, "", "   ")
    case '$':
        message := new(tq.Message)
        message.Decode(data)
        out, _ = json.MarshalIndent(message, "", "   ")
    case 0x00:
        message := new(cmpp.Message)
        message.Decode(data)
        out, _ = json.MarshalIndent(message, "", "   ")
    }

    if strings.HasPrefix(hexStr, "Gpslocation") {
        message := new(th.Message)
        message.Decode(data)
        out, _ = json.MarshalIndent(message, "", "   ")
    } else if strings.HasPrefix(hexStr, "Lbslocation") {
        message := new(th.Message)
        message.Decode(data)
        out, _ = json.MarshalIndent(message, "", "   ")
    } else if strings.HasPrefix(hexStr, "#006") {
        message := new(th.Message)
        message.Decode(data)
        out, _ = json.MarshalIndent(message, "", "   ")
    } else if strings.HasPrefix(hexStr, "#007") {
        message := new(th.Message)
        message.Decode(data)
        out, _ = json.MarshalIndent(message, "", "   ")
    }

    fmt.Println(string(out))
}
