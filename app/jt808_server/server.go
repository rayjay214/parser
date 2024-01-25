package main

import (
    "fmt"
    "github.com/rayjay214/parser/jt808"
    "github.com/rayjay214/parser/jt808/extra"
    "github.com/rayjay214/parser/server"
    log "github.com/sirupsen/logrus"
)

func handle0100(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0100)
    log.Info(entity)
    session.ReplyRegister(message)
}

func handle0102(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0102)
    log.Info(entity)

    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func handle0002(session *server.Session, message *jt808.Message) {
    entity := message.Body.(*jt808.T808_0x0002)
    log.Info(entity)

    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

// 处理上报位置
func handle0200(session *server.Session, message *jt808.Message) {
    // 打印消息
    entity := message.Body.(*jt808.T808_0x0200)
    fields := log.Fields{
        "Imei": message.Header.Imei,
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

    // 回复平台应答
    session.Reply(message, jt808.T808_0x8100_ResultSuccess)
}

func main() {
    server, _ := server.NewServer(server.Options{
        Keepalive:       60,
        AutoMergePacket: true,
        CloseHandler:    nil,
    })
    server.AddHandler(jt808.MsgT808_0x0100, handle0100)
    server.AddHandler(jt808.MsgT808_0x0102, handle0102)
    server.AddHandler(jt808.MsgT808_0x0002, handle0002)
    server.AddHandler(jt808.MsgT808_0x0200, handle0200)
    server.Run("tcp", 8808)
}
