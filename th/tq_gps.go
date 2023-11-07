package th

import (
    _ "fmt"
    "github.com/rayjay214/parser/common"
    _ "strconv"
    "strings"
    "time"
)

type TH_GPS struct {
    Proto       string
    Imei        string
    Time        time.Time
    Lat         string
    Lng         string
    Speed       string
    Direction   string
    Altitude    string
    Temperature string
    Humility    string
    Status      string
    Battery     string
    GSM         string
    GPSSignal   string
    OnOffStatus string
    Flag        string
}

func (entity *TH_GPS) MsgID() MsgID {
    return Msg_GPS
}

func (entity *TH_GPS) Encode() ([]byte, error) {
    writer := common.NewWriter()

    //todo

    return writer.Bytes(), nil
}

func (entity *TH_GPS) Decode(data []byte) (int, error) {
    strData := string(data)
    strList := strings.Split(strData, "|")
    entity.Proto = strList[0]
    entity.Imei = strList[1]
    strTime := strList[2]
    strLngLat := strList[3]
    entity.Speed = strList[4]
    entity.Direction = strList[5]
    entity.Altitude = strList[6]
    entity.Temperature = strList[7]
    entity.Humility = strList[8]
    entity.Status = strList[9]
    entity.Battery = strList[10]
    entity.GSM = strList[11]
    entity.GPSSignal = strList[12]
    entity.OnOffStatus = strList[13]
    entity.Flag = strList[14]

    entity.Lng = strings.Split(strLngLat, ",")[0]
    entity.Lat = strings.Split(strLngLat, ",")[1]

    entity.Time, _ = time.ParseInLocation("060102150405", strTime, time.UTC)

    return len(data), nil
}
