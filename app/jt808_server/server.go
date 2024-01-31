package main

import (
    "github.com/rayjay214/parser/jt808"
    "github.com/rayjay214/parser/server"
    "github.com/rayjay214/parser/service"
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

    go service.StartRpc(server)
    server.Run("tcp", 8808)
}
