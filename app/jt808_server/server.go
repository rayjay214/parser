package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rayjay214/parser/app/jt808_server/service"
	"github.com/rayjay214/parser/jt808"
	"github.com/rayjay214/parser/server"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"time"
)

func init() {
	path := "log/server.log"
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(60*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(7*24)*time.Hour),
	)
	log.SetOutput(writer)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	log.Info("init log done")
}

func main() {
	server, _ := server.NewServer(server.Options{
		Keepalive:       420,
		AutoMergePacket: true,
		CloseHandler:    nil,
	})
	server.AddHandler(jt808.MsgT808_0x0100, handle0100)
	server.AddHandler(jt808.MsgT808_0x0102, handle0102)
	server.AddHandler(jt808.MsgT808_0x0002, handle0002)
	server.AddHandler(jt808.MsgT808_0x0200, handle0200)
	server.AddHandler(jt808.MsgT808_0x0704, handle0704)
	server.AddHandler(jt808.MsgT808_0x0808, handle0808)
	server.AddHandler(jt808.MsgT808_0x1007, handle1007)
	server.AddHandler(jt808.MsgT808_0x1107, handle1107)
	server.AddHandler(jt808.MsgT808_0x1300, handle1300)
	server.AddHandler(jt808.MsgT808_0x0116, handle0116)
	server.AddHandler(jt808.MsgT808_0x0117, handle0117)
	server.AddHandler(jt808.MsgT808_0x0109, handle0109)
	server.AddHandler(jt808.MsgT808_0x0003, handle0003)
	server.AddHandler(jt808.MsgT808_0x0105, handle0105)
	server.AddHandler(jt808.MsgT808_0x0108, handle0108)
	server.AddHandler(jt808.MsgT808_0x0210, handle0210)
	server.AddHandler(jt808.MsgT808_0x0115, handle0115)
	server.AddHandler(jt808.MsgT808_0x0120, handle0120)
	server.AddHandler(jt808.MsgT808_0x0118, handle0118)
	server.AddHandler(jt808.MsgT808_0x0119, handle0119)
	server.AddHandler(jt808.MsgT808_0x0001, handle0001)
	server.AddHandler(jt808.MsgT808_0x0107, handle0107)
	server.AddHandler(jt808.MsgT808_0x0112, handle0112)
	server.AddHandler(jt808.MsgT808_0x1006, handle1006)

	storage.InitCass("47.107.69.24")
	storage.InitMinio("114.215.190.173:9000")
	storage.InitRedis("47.107.69.24:6480")
	storage.InitMysql("47.107.69.24:8000")

	go service.StartRpc(server)
	server.Run("tcp", 8808)
}
