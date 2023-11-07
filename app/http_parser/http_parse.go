package main

import (
    "encoding/binary"
    "encoding/hex"
    "encoding/json"
    _ "fmt"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "net/http"
    "os"
    "parser/cmpp"
    "parser/common"
    "parser/ipc"
    "parser/jt808"
    "parser/kks"
    "parser/ota"
    "parser/th"
    "parser/tq"
    "strings"
)

type RawInfo struct {
    Raws []string `json:"raws"`
}

type ParsedInfo struct {
    Parsed []string `json:"parsed"`
}

func init() {
    f, err := os.OpenFile("parser.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
    if err != nil {
        log.Fatal(err)
    }
    log.SetOutput(f)
    log.SetLevel(log.InfoLevel)
    log.SetReportCaller(true)
    log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
    log.Info("init log done")
}

func ParseHandler(writer http.ResponseWriter, request *http.Request) {
    v := request.URL.Query()
    hexStr := v["hexstr"][0]

    if len(hexStr) == 0 {
        return
    }

    var err interface{}
    var data []byte
    var out []byte
    if hexStr[0] == '*' {
        data = []byte(hexStr)
    } else if strings.HasPrefix(hexStr, "Gpslocation") || strings.HasPrefix(hexStr, "Lbslocation") || strings.HasPrefix(hexStr, "#006") || strings.HasPrefix(hexStr, "#007") {
        data = []byte(hexStr)
        log.Warn("Parse gps location: ")
        message := new(th.Message)
        message.Decode(data)
        out, _ = json.MarshalIndent(message, "", "   ")
        writer.Write(out)
        return
    } else {
        data, err = hex.DecodeString(hexStr)
        if err != nil {
            return
        }
    }

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
        message.Decode(transformed_data)
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
    case 0x86:
        message := new(ipc.Message)
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
    writer.Write(out)
}

func BatchParseHandler(writer http.ResponseWriter, request *http.Request) {
    var parsedInfo ParsedInfo
    body, _ := ioutil.ReadAll(request.Body)
    var info RawInfo
    json.Unmarshal(body, &info)

    var data []byte
    var err interface{}
    parsedInfo.Parsed = make([]string, 0)
    for i := range info.Raws {
        hexStr := info.Raws[i]

        if len(hexStr) == 0 {
            log.Warn("Empty str")
            parsedInfo.Parsed = append(parsedInfo.Parsed, "")
            continue
        }

        var out []byte
        if hexStr[0] == '*' {
            data = []byte(hexStr)
        } else if strings.HasPrefix(hexStr, "Gpslocation") || strings.HasPrefix(hexStr, "Lbslocation") || strings.HasPrefix(hexStr, "#006") || strings.HasPrefix(hexStr, "#007") {
            data = []byte(hexStr)
            message := new(th.Message)
            message.Decode(data)
            out, _ = json.MarshalIndent(message, "", "   ")
            parsedInfo.Parsed = append(parsedInfo.Parsed, string(out))
            continue
        } else {
            data, err = hex.DecodeString(hexStr)
            if err != nil {
                log.Warn("Decode Error: ", err)
                parsedInfo.Parsed = append(parsedInfo.Parsed, "parse error")
                continue
            }
        }

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
            message.Decode(transformed_data)
            out, _ = json.Marshal(message)
        case 0x78:
            message := new(kks.Message_0x78)
            message.Decode(data)
            out, _ = json.Marshal(message)
        case 0x79:
            message := new(kks.Message_0x79)
            message.Decode(data)
            out, _ = json.Marshal(message)
        case 0x68:
            message := new(ota.Message)
            message.Decode(data)
            out, _ = json.Marshal(message)
        case '*':
            message := new(tq.Message)
            message.Decode(data)
            out, _ = json.Marshal(message)
        case '$':
            message := new(tq.Message)
            message.Decode(data)
            out, _ = json.Marshal(message)
        }
        parsedInfo.Parsed = append(parsedInfo.Parsed, string(out))
    }
    p, _ := json.Marshal(parsedInfo)

    writer.Write([]byte(p))
}

func main() {
    http.HandleFunc("/parse", ParseHandler)
    http.HandleFunc("/batch_parse", BatchParseHandler)
    http.ListenAndServe("0.0.0.0:8081", nil)
}
