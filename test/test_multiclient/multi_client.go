package main

import (
	"encoding/hex"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/rayjay214/parser/protocol/jt808"
	"github.com/shopspring/decimal"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var g_seqno uint32 = 0

type location struct {
	lat       float64
	lon       float64
	speed     uint16
	direction uint16
}

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

func initCassSession() (session *gocql.Session) {
	fmt.Printf("begin init")
	cluster := gocql.NewCluster("47.104.8.119") //replace PublicIP with the IP addresses used by your cluster.
	cluster.Consistency = gocql.LocalOne
	cluster.ProtoVersion = 3
	cluster.ConnectTimeout = time.Second * 10
	cluster.Keyspace = "slxk"
	session, err := cluster.CreateSession()
	if err != nil {
		return nil
	}
	fmt.Printf("init over")
	return session
}

func sender(conn *net.TCPConn, wg *sync.WaitGroup, s_location []location, i int) {
	defer wg.Done()
	//login
	var imei uint64 = 77777000000 + uint64(i)
	authkey := strconv.FormatUint(imei, 10)
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
	fmt.Printf("%s", hex.Dump(data))

	if _, err = conn.Write(data); err != nil {
		fmt.Printf("write failed , err : %v\n", err)
		return
	}

	//send gps packet recursively
	idx := 0
	for {
		if idx >= len(s_location) {
			idx = 0
		}

		loc_message := jt808.Message{
			Header: jt808.Header{
				Imei:        imei,
				MsgSerialNo: nextID(),
			},
			Body: &jt808.T808_0x0200{
				Alarm:     512,
				Status:    2,
				Lat:       decimal.NewFromFloat(s_location[idx].lat),
				Lng:       decimal.NewFromFloat(s_location[idx].lon),
				Altitude:  0,
				Speed:     s_location[idx].speed,
				Direction: s_location[idx].direction,
				Time:      time.Unix(time.Now().Unix(), 0),
			},
		}
		fmt.Println(loc_message)
		data, err := loc_message.Encode()
		fmt.Printf("%s", hex.Dump(data))

		if _, err = conn.Write(data); err != nil {
			fmt.Printf("write failed , err : %v\n", err)
			break
		}

		time.Sleep(time.Duration(30) * time.Second)
		idx += 1
	}
}

func receiver(conn *net.TCPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
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
	}
}

func main() {
	var wg sync.WaitGroup
	session := initCassSession()

	//s_location := make([]location, 1000)
	var s_location []location
	var ilat, ilng uint32
	var speed, direction uint16
	iter := session.Query("SELECT flon, flat, fspeed, fdirection FROM tkv_location where fimei=99999999035 and fdate=20210122").Iter()
	for iter.Scan(&ilng, &ilat, &speed, &direction) {
		loc := location{
			float64(ilat) / 1000000,
			float64(ilng) / 1000000,
			speed,
			direction}
		fmt.Println(loc)
		s_location = append(s_location, loc)
	}

	for i := 1; i < 10000; i++ {
		tcpAddr, err := net.ResolveTCPAddr("tcp", "47.104.8.119:8881")
		if err != nil {
			panic(err)
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Printf("connect failed, err : %v\n", err.Error())
			return
		}

		defer conn.Close()
		wg.Add(1)
		go sender(conn, &wg, s_location, i)
		wg.Add(1)
		go receiver(conn, &wg)
		time.Sleep(time.Duration(10) * time.Millisecond)
	}

	wg.Wait()
}
