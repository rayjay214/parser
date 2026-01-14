package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rayjay214/parser/app/hl3g_server/service"
	"github.com/rayjay214/parser/protocol/hl3g"
	"github.com/rayjay214/parser/server_base/hl3g_base"
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
	server, _ := hl3g_base.NewServer(hl3g_base.Options{
		Keepalive:       420,
		AutoMergePacket: true,
		CloseHandler:    nil,
	})
	server.AddHandler(hl3g.Msg_LK2, handleLK2)
	server.AddHandler(hl3g.Msg_CCID, handleCCID)
	server.AddHandler(hl3g.Msg_GS1, handleGS1)
	server.AddHandler(hl3g.Msg_UD, handleUD)
	server.AddHandler(hl3g.Msg_UD2, handleUD2)
	server.AddHandler(hl3g.Msg_AL, handleAL)

	storage.InitCass(storage.Conf.Cassandra.Host)
	storage.InitMinio(storage.Conf.Minio.Host)
	storage.InitRedis(storage.Conf.Redis.Host)
	storage.InitMysql(storage.Conf.Mysql.Host)

	go service.StartRpc(server)
	server.Run("tcp", storage.Conf.Jt808Server.Port)
}
