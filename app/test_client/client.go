package main

import (
	//"encoding/hex"
	"fmt"
	"net"
	"parser/jt808"
	"sync/atomic"
	"time"
)

var g_seqno uint32 = 0

func nextID() uint16 {
    var id uint32
    for {
        id = atomic.LoadUint32(&g_seqno)
        if id == 0xff {
            if atomic.CompareAndSwapUint32(&g_seqno, id, 1) {
                id = 1
                break
            }
        } else if atomic.CompareAndSwapUint32(&g_seqno, id, id+1) {
            id += 1
            break
        }
    }
    return uint16(id)
}

func main() {

    tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8881")
    if err != nil {
        panic(err)
    }

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        fmt.Printf("connect failed, err : %v\n", err.Error())
        return
    }

    defer conn.Close()

    for {
        // 终端鉴权
        var imei uint64 = 12345123451
        var authkey = "12345123451"
        message := jt808.Message{
            Header: jt808.Header{
                Imei:        imei,
                MsgSerialNo: nextID(),
            },
            Body: &jt808.T808_0x0102{
                AuthKey: authkey,
            },
        }
        data, err := message.Encode()

        if _, err = conn.Write(data); err != nil {
            fmt.Printf("write failed , err : %v\n", err)
            break
        }

        var p jt808.Protocol
        codec, err := p.NewCodec(conn)
        if err != nil {
            panic(err)
        }
        msg, err := codec.Receive()
        resp := msg.(jt808.Message)

        if resp.Header.MsgID == jt808.MsgT808_0x8001 {
            body := resp.Body.(*jt808.T808_0x8001)
            fmt.Println(body.Result)
        }

        time.Sleep(time.Duration(10) * time.Second)
    }
}
