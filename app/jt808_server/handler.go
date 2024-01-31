package main

import (
    "fmt"
    "github.com/rayjay214/parser/common"
    "github.com/rayjay214/parser/jt808"
    "github.com/rayjay214/parser/jt808/extra"
    "github.com/rayjay214/parser/server"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
)

func handle0100(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0100)
    log.Infof("handle 0100 %v", entity)

    session.ReplyRegister(message)
}

func handle0102(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0102)
    log.Infof("handle 0102 %v", entity)

    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0002(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0002)
    log.Infof("handle 0002 %v", entity)

    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0808(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0808)
    log.Infof("handle 0808 %v", entity)
}

// 处理上报位置
func handle0200(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0200)
    log.Infof("handle 0200 %v", entity)
    handleLocation(message.Header.Imei, entity)

    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0704(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0704)
    log.Infof("handle 0704 %v", entity)

    for _, item := range entity.Items {
        handleLocation(message.Header.Imei, &item)
    }

    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handleLocation(imei uint64, entity *jt808.T808_0x0200) {
    fields := log.Fields{
        "Imei": imei,
        "警告": fmt.Sprintf("0x%x", entity.Alarm),
        "状态": fmt.Sprintf("0x%x", entity.Status),
        "纬度": entity.Lat,
        "经度": entity.Lng,
        "海拔": entity.Altitude,
        "速度": entity.Speed,
        "方向": entity.Direction,
        "时间": entity.Time,
    }

    for _, ext := range entity.Extras {
        switch ext.ID() {
        case extra.Extra_0x01{}.ID():
            fields["行驶里程"] = ext.(*extra.Extra_0x01).Value()
        case extra.Extra_0x02{}.ID():
            fields["剩余油量"] = ext.(*extra.Extra_0x02).Value()
        }
    }
    log.WithFields(fields).Info("上报终端位置信息")
}

func handle1007(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x1007)
    log.Infof("handle 1007 %v", entity)

    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle1107(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x1107)
    log.Infof("handle 1107 %v", entity)

    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle1300(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x1300)
    log.Infof("handle 1300 %v", entity)
}

func handle0116(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0116)
    session.UserData["schedule"] = 0
    session.UserData["short_record_writer"] = common.NewWriter()
    log.Infof("handle 0116 %v", entity)
}

func handle0117(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0117)
    session.UserData["schedule"] = float32(entity.PkgNo) * 100.0 / float32(entity.PkgSize)
    log.Infof("schedule is %v", session.UserData["schedule"])
    writer, ok := session.UserData["short_record_writer"].(common.Writer)
    if ok {
        log.Infof("append pkt %v", entity.PkgNo)
        buffer, _ := ioutil.ReadAll(entity.Packet)
        writer.Write(buffer)
        log.Infof("buffer len %v", len(writer.Bytes()))
    }

    if entity.PkgNo == entity.PkgSize {
        file, _ := os.Create("aa.amr")
        defer file.Close()
        leng, err := file.Write(writer.Bytes())
        log.Infof("%v:%v", leng, err)
    }

    data, _ := message.Encode()
    log.Infof("0117 raw msg %x", common.GetHex(data))
    session.ReplyShortRecord(entity.PkgNo)
}

func handle0109(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0109)
    session.ReplyTime()
    log.Infof("handle 0109 %v", entity)
}
