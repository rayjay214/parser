package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rayjay214/parser/app/gt06_server/service"
	"github.com/rayjay214/parser/protocol/gt06"
	"github.com/rayjay214/parser/server_base/gt06_base"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"time"
)

func init() {
	path := "log/server.log"
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(5*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	log.SetOutput(writer)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	log.Info("init log done")
}

func main() {
	storage.LoadConfig("config.ini")
	storage.InitRawLog()
	server, _ := gt06_base.NewServer(gt06_base.Options{
		Keepalive:       420,
		AutoMergePacket: true,
		CloseHandler:    nil,
	})
	server.AddHandler(gt06.Msg_0x01, handle01)
	server.AddHandler(gt06.Msg_0x12, handle12)
	server.AddHandler(gt06.Msg_0xa1, handleA1)
	server.AddHandler(gt06.Msg_0x94, handle94)
	server.AddHandler(gt06.Msg_0x20, handle20)
	server.AddHandler(gt06.Msg_0x13, handle13)
	server.AddHandler(gt06.Msg_0x16, handle16)

	storage.InitCass(storage.Conf.Cassandra.Host)
	storage.InitMinio(storage.Conf.Minio.Host)
	storage.InitRedis(storage.Conf.Redis.Host)
	storage.InitMysql(storage.Conf.Mysql.Host)

	go service.StartRpc(server)
	server.Run("tcp", storage.Conf.Jt808Server.Port)
}
