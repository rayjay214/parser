package main

import (
    "fmt"
    "github.com/rayjay214/parser/common"
    "github.com/rayjay214/parser/jt808"
    "github.com/rayjay214/parser/jt808/extra"
    "github.com/rayjay214/parser/server"
    "github.com/rayjay214/parser/storage"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
    "time"
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
    session.Protocol = 2 //2013
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

    cassSession, err := storage.GetSession()
    defer cassSession.Close()
    if err != nil {
        log.Warnf("get cassandra failed %v", err)
    }
    err = insertLocation(cassSession, entity, imei)
    if err != nil {
        log.Warnf("insert location err %v", err)
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
    session.UserData["short_record"] = ShortRecord{
        Imei:      session.ID(),
        Writer:    common.NewWriter(),
        StartTime: time.Now(),
        Schedule:  0.0,
    }
    log.Infof("handle 0116 %v", entity)
}

func handle0117(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0117)

    /*
       这里Encode的时候，会把packet的内容写走
       log.Infof("body is %v", entity)
       d, _ := message.Encode()
       log.Infof("0117 raw msg %x", common.GetHex(d))
       bdata, _ := entity.Encode()
       log.Infof("0117 raw body %x", common.GetHex(bdata))
    */

    shortRecord, ok := session.UserData["short_record"].(ShortRecord)
    if ok {
        shortRecord.Schedule = float32(entity.PkgNo) * 100.0 / float32(entity.PkgSize)
        log.Infof("schedule is %.2f", shortRecord.Schedule)
        buffer, _ := ioutil.ReadAll(entity.Packet)
        shortRecord.Writer.Write(buffer)
    }

    if entity.PkgNo == entity.PkgSize {
        fileName := fmt.Sprintf("record/%v_%v.amr", shortRecord.Imei, shortRecord.StartTime.Unix())
        file, _ := os.Create(fileName)
        defer file.Close()
        file.Write(shortRecord.Writer.Bytes())
        shortRecord.Writer.Reset()
    }

    data, _ := message.Encode()
    log.Infof("0117 raw msg %x", common.GetHex(data))
    session.ReplyShortRecord(entity.PkgNo)
}

func handle0109(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0109)
    log.Infof("handle 0109 %v", entity)
    session.ReplyTime()
}

func handle0003(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0003)
    log.Infof("handle 0109 %v", entity)
    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0105(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0105)
    session.Reply8125()
    log.Infof("handle 0105 %v", entity)
}

func handle0108(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0108)
    log.Infof("handle 0108 %v", entity)
    session.Reply8108()
}

func handle0210(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0210)
    log.Infof("handle 0210 %v", entity)
    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0115(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0115)
    log.Infof("handle 0115 %v", entity)
    delete(session.UserData, "short_record")
    //session.Reply8115(entity.SessionId)
}

func handle0120(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0120)
    log.Infof("handle 0120 %v", entity)
    session.UserData["vor_record"] = &VorRecord{
        Imei:        session.ID(),
        Writer:      common.NewWriter(),
        StartTime:   time.Now(),
        EndTime:     time.Now(),
        FirstPacket: true,
        PkgCnt:      0,
    }
    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0118(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0118)

    currBeginTime := entity.Time

    vorRecord, ok := session.UserData["vor_record"].(*VorRecord)
    if !ok {
        session.UserData["vor_record"] = &VorRecord{
            Imei:        session.ID(),
            Writer:      common.NewWriter(),
            StartTime:   time.Now(),
            EndTime:     time.Now(),
            FirstPacket: true,
            PkgCnt:      0,
        }
        vorRecord, _ = session.UserData["vor_record"].(*VorRecord)
    }

    if vorRecord.FirstPacket {
        vorRecord.StartTime = entity.Time
        vorRecord.EndTime = vorRecord.StartTime.Add(time.Second * 10)
        vorRecord.FirstPacket = false
    }
    if entity.PkgNo == 1 { //组包
        if currBeginTime.Sub(vorRecord.EndTime).Seconds() < 3 &&
            currBeginTime.Sub(vorRecord.StartTime).Seconds() < 59 &&
            vorRecord.PkgCnt < 47 {
            vorRecord.EndTime = currBeginTime.Add(time.Second * 10)
            buffer, err := ioutil.ReadAll(entity.Packet)
            log.Infof("rayjay buffer len %v, err is %v", len(buffer), err)
            vorRecord.Writer.Write(buffer)
            vorRecord.PkgCnt += 1
        } else { //上报当前缓存录音
            fileName := fmt.Sprintf("vrecord/%v_%v.amr", vorRecord.Imei, vorRecord.StartTime.Unix())
            log.Infof("upload record %v", fileName)
            file, _ := os.Create(fileName)
            defer file.Close()
            file.Write(vorRecord.Writer.Bytes())
            vorRecord.Writer.Reset()
            //重新初始化并写入第一个包
            vorRecord.PkgCnt = 0
            vorRecord.FirstPacket = true
            buffer, err := ioutil.ReadAll(entity.Packet)
            log.Infof("rayjay buffer len %v, err is %v", len(buffer), err)
            vorRecord.Writer.Write(buffer)
            vorRecord.PkgCnt += 1
        }
    } else {
        buffer, err := ioutil.ReadAll(entity.Packet)
        log.Infof("rayjay buffer len %v, err is %v", len(buffer), err)
        vorRecord.Writer.Write(buffer)
        vorRecord.PkgCnt += 1
    }
    log.Infof("vor %v", *vorRecord)

    session.ReplyVorRecord(entity)
}

func handle0119(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0119)
    log.Infof("handle 0119 %v", entity)
    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0001(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0001)
    log.Infof("handle 0001 %v", entity)
}
