package main

import (
	"github.com/rayjay214/parser/app/jt808_server/service"
	"github.com/rayjay214/parser/jt808"
	"github.com/rayjay214/parser/server"
	"github.com/rayjay214/parser/storage"
)

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

	storage.InitCass("47.107.69.24")
	storage.InitMinio("114.215.190.173:9000")
	storage.InitRedis("47.107.69.24:6480")

	go service.StartRpc(server)
	server.Run("tcp", 8808)
}
