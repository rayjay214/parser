package main

import (
    //"encoding/hex"
    "fmt"
    "github.com/rayjay214/parser/ipc"
    "net"
    "time"
)

func sendMsg(conn *net.TCPConn) {
    var msg ipc.Message
    msg.Header.Prefix = 0x8686
    msg.Header.MsgId = ipc.Msg_0x0003
    msg.Header.Seq = 11
    msg.Header.UidLen = 15
    msg.Header.Uid = "123456789111111"
    msg.Body = new(ipc.Body_0x0003)
    body0003 := msg.Body.(*ipc.Body_0x0003)
    body0003.Ip = "192.168.1.2"
    body0003.Port = 1234

    data, _ := msg.Encode()
    fmt.Printf("hex %x\n", ipc.GetHex(data))
    conn.Write(data)

    time.Sleep(1 * time.Second)
    buf := make([]byte, 128)
    _, err := conn.Read(buf[:])
    if err == nil {
        fmt.Printf(string(buf))
    }
}

func main() {
    tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8882")
    if err != nil {
        panic(err)
    }

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        fmt.Printf("connect failed, err : %v\n", err.Error())
        return
    }

    sendMsg(conn)

    defer conn.Close()

}
