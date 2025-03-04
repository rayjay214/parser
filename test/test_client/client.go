package main

import (
	//"encoding/hex"
	"fmt"
	"github.com/rayjay214/parser/protocol/jt808"
	"github.com/shopspring/decimal"
	"net"
	"sync/atomic"
	"time"
)

var g_seqno uint32 = 2

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

	tcpAddr, err := net.ResolveTCPAddr("tcp", "114.215.190.173:8881")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}

	defer conn.Close()

	// 终端鉴权
	var imei uint64 = 13320774980
	var authkey = "13320774980"
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
	fmt.Printf("%x\n", data)
	if _, err = conn.Write(data); err != nil {
		fmt.Printf("write failed , err : %v\n", err)
		return
	}

	fmt.Println("send ok")

	locMessage := jt808.Message{
		Header: jt808.Header{
			Imei:        imei,
			MsgSerialNo: nextID(),
		},
		Body: &jt808.T808_0x0200{
			Alarm:     0,
			Status:    2,
			Lat:       decimal.NewFromFloat(22.677786),
			Lng:       decimal.NewFromFloat(114.145426),
			Altitude:  0,
			Speed:     0,
			Direction: 0,
			Time:      time.Unix(time.Now().Unix(), 0),
		},
	}

	data, err = locMessage.Encode()
	if _, err = conn.Write(data); err != nil {
		fmt.Printf("write failed , err : %v\n", err)
		return
	}

	time.Sleep(time.Duration(10) * time.Second)

}
