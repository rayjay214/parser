package main

import (
    "bytes"
    _ "encoding/binary"
    "encoding/hex"
    "errors"
    "flag"
    _ "fmt"
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    "github.com/rayjay214/parser/protocol/jt808"
    log "github.com/sirupsen/logrus"
    "io"
    "net"
    "os"
    "strconv"
    "syscall"
    "time"
)

func init() {
	path := "log/server.log"
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(15*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(7*24)*time.Hour),
	)
	log.SetOutput(writer)
	log.SetLevel(log.DebugLevel)
	//log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	log.Info("init log done")
}

var g_device map[uint64]time.Time
var g_err_device map[uint64]time.Time

//routine cnt
var g_connection_routine int
var g_gwclient_routine int
var g_gwreceiver_routine int

func dump() {
	var data string
	now := time.Now()
	for k, v := range g_device {
		strImei := strconv.Itoa(int(k))
		if (now.Unix() - v.Unix()) > 7*60 {
			log.Info(strImei, " offline")
			delete(g_device, k)
			continue
		}
		data += strImei + ":" + v.Format("2006-01-02 15:04:05") + "\n"
	}

	file, err := os.OpenFile("device.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Error("open file failed", err)
		return
	}
	defer file.Close()
	file.WriteString(data)
}

func dump_err() {
	var data string
	now := time.Now()
	for k, v := range g_err_device {
		strImei := strconv.Itoa(int(k))
		if (now.Unix() - v.Unix()) > 7*60 {
			log.Info(strImei, " offline or becoming normal")
			delete(g_err_device, k)
			continue
		}
		data += strImei + ":" + v.Format("2006-01-02 15:04:05") + "\n"
	}

	file, err := os.OpenFile("err_device.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Error("open file failed", err)
		return
	}
	defer file.Close()
	file.WriteString(data)
}

func dump_routine_cnt() {
	var data string

	data += "connection routine:" + strconv.Itoa(g_connection_routine) + "\n"
	data += "gwclient routine:" + strconv.Itoa(g_gwclient_routine) + "\n"
	data += "gwreceiver routine:" + strconv.Itoa(g_gwreceiver_routine) + "\n"

	file, err := os.OpenFile("routine_cnt.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Error("open file failed", err)
		return
	}
	defer file.Close()
	file.WriteString(data)
}

func addCacheTicker() {
	ticker := time.NewTicker(time.Second * 120)
	go func() {
		for {
			<-ticker.C
			dump()
			dump_err()
			dump_routine_cnt()
		}
	}()
}

func main() {
	var listenAddr string
	var transAddr string
	flag.StringVar(&listenAddr, "listenAddr", "0.0.0.0:8881", "监听地址")
	flag.StringVar(&transAddr, "transAddr", "47.104.8.119:8881", "转发GW地址")

	flag.Parse()

	/*
	       ch := make(chan int)
	       close(ch)
	       ch <- 10
	   	logFile, err := os.OpenFile(listenAddr + ".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	       if err != nil {
	           fmt.Println("open log file failed, err:", err)
	           return
	       }
	       log.SetOutput(logFile)
	       log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
	*/

	g_device = make(map[uint64]time.Time)
	g_err_device = make(map[uint64]time.Time)
	addCacheTicker()

	tcpAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Error("Listen Error: ", err)
		return
	}
	for {
		conn, err := ln.Accept()
		tcpConn := conn.(*net.TCPConn)
		if err != nil {
			log.Error("Accept Error: ", err)
			continue
		}
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(time.Second * 60)
		go handleConnection(tcpConn, transAddr)
	}

}

func transform(data []byte) []byte {
	var buffer bytes.Buffer
	buffer.WriteByte(0x7e)
	for _, b := range data[1 : len(data)-1] {
		if b == 0x7e {
			buffer.WriteByte(0x7d)
			buffer.WriteByte(0x02)
		} else if b == 0x7d {
			buffer.WriteByte(0x7d)
			buffer.WriteByte(0x01)
		} else {
			buffer.WriteByte(b)
		}
	}
	buffer.WriteByte(0x7e)
	return buffer.Bytes()
}

func deTransform(databuf []byte, pkglen *int, origlen *int, t_buf []byte) int {
	if len(databuf) == 0 {
		return -1
	}

	if databuf[0] != 0x7e || *pkglen < 2 {
		return -2
	}

	var end int
	for end = 1; end < len(databuf); end++ {
		if databuf[end] == 0x7e {
			break
		}
	}
	if end < *pkglen { //找到一个完整的包
		*pkglen = end + 1
		*origlen = end + 1
	} else {
		return 9999
	}

	curr := 1
	skipped := 0

	t_buf[0] = databuf[0]

	for i := 1; i < *pkglen-1; i++ {
		if databuf[curr] == 0x7d {
			if i+1 >= *pkglen-1 {
				return -3
			}
			if databuf[curr+1] == 0x02 {
				t_buf[curr-skipped] = 0x7e
			} else if databuf[curr+1] == 0x01 {
				t_buf[curr-skipped] = 0x7d
			} else {
				return -4
			}
			i++
			curr++
			skipped++
		} else {
			t_buf[curr-skipped] = databuf[curr]
		}
		curr++
	}

	t_buf[curr-skipped] = databuf[curr]
	*pkglen = *pkglen - skipped

	return 0
}

func handleConnection(conn *net.TCPConn, transAddr string) {
	g_connection_routine++
	defer conn.Close()

	log.Info("Client: ", conn.RemoteAddr(), " Connected")

	msgbuf := bytes.NewBuffer(make([]byte, 0, 10240))

	databuf := make([]byte, 4096)

	recv := 0

	//控制协程退出标志

	//数据传输
	ch := make(chan []byte, 10)

	//控制协程退出
	ctrl_ch := make(chan int, 1)

	go client(ch, transAddr, ctrl_ch, conn)

	var m_imei uint64

	for {
		//  Read the data
		err := conn.SetReadDeadline(time.Now().Add(300 * time.Second))
		if err != nil {
			log.Error("SetReadDeadline failed:", err)
			ctrl_ch <- 1
			close(ch)
			g_connection_routine--
			return
		}

		n, err := conn.Read(databuf)
		if err != nil {
			if err == io.EOF {
				log.Warn("Client ", conn.RemoteAddr(), " exit:")
			} else if errors.Is(err, syscall.ECONNRESET) {
				g_err_device[m_imei] = time.Now()
				log.Warn(m_imei, " Read error: ", err)
			} else if errors.Is(err, os.ErrDeadlineExceeded) {
				log.Warn(m_imei, " Read error: ", err)
			} else {
				log.Warn(m_imei, " Read error: ", err)
			}
			ctrl_ch <- 1
			close(ch)
			g_connection_routine--
			return
		}

		//  Data is added to the message buffer
		n, err = msgbuf.Write(databuf[:n])
		if err != nil {
			log.Error("Buffer write error: ", err)
			ctrl_ch <- 1
			close(ch)
			g_connection_routine--
			return
		}

		//  Message segmentation loop
		var ret int
		for {
			tmpdata := msgbuf.Bytes()
			if len(tmpdata) <= 0 {
				break
			}

			var pkglen int
			var origlen int
			pkglen = len(tmpdata)
			var t_buf []byte = make([]byte, 10240)
			ret = deTransform(tmpdata, &pkglen, &origlen, t_buf)
			if ret < 0 {
				log.Error("detransform failed ", ret)
				ctrl_ch <- 1
				close(ch)
				g_connection_routine--
				return
			} else if ret == 9999 {
				log.Info("package incomplete")
				break
			}
			log.Debug("receive whole msg:", hex.EncodeToString(t_buf[:pkglen]))
			ret = handleMsg(conn, t_buf[:pkglen], &m_imei)
			if ret != 0 {
				log.Warn("handlemsg failed")
				ctrl_ch <- 1
				close(ch)
				g_connection_routine--
				return
			}
			ch <- t_buf[:pkglen]

			msgbuf.Next(origlen)
			recv += 1
		}
	}
	g_connection_routine--
	close(ch)
}

var g_seq uint16

func handleMsg(conn *net.TCPConn, databuf []byte, m_imei *uint64) int {
	//response
	var err error
	var respdata []byte
	message := new(jt808.Message)
	transedData := transform(databuf)
	err = message.Decode(transedData)
	if err != nil {
		log.Error("decode message failed ", err)
		return -1
	}
	var imei uint64 = message.Header.Imei
	*m_imei = imei
	//now := time.Now().Format("2006-01-02 15:04:05")
	g_device[imei] = time.Now()

	/*
		_, ok := g_err_device[imei]
		if ok {
			log.Printf("handle errdevice %v", imei)
			return -1
		}
	*/

	switch {
	case int(message.Header.MsgID) == 0x0100:
		msg8100 := jt808.Message{
			Header: jt808.Header{
				Imei:        imei,
				MsgSerialNo: g_seq,
			},
			Body: &jt808.T808_0x8100{
				MsgSerialNo: message.Header.MsgSerialNo,
				Result:      0,
				AuthKey:     strconv.Itoa(int(imei)),
			},
		}
		respdata, err = msg8100.Encode()
	case int(message.Header.MsgID) == 0x0109:
		now := time.Now()

		msg8109 := jt808.Message{
			Header: jt808.Header{
				Imei:        imei,
				MsgSerialNo: g_seq,
			},
			Body: &jt808.T808_0x8109{
				Year:   uint16(now.Year()),
				Month:  byte(now.Month()),
				Day:    byte(now.Day()),
				Hour:   byte(now.Hour()),
				Minute: byte(now.Minute()),
				Second: byte(now.Second()),
				Result: 0,
			},
		}
		respdata, err = msg8109.Encode()
	case int(message.Header.MsgID) == 0x0808:
		respdata = nil
	default:
		msg8001 := jt808.Message{
			Header: jt808.Header{
				Imei:        imei,
				MsgSerialNo: g_seq,
			},
			Body: &jt808.T808_0x8001{
				ReplyMsgSerialNo: message.Header.MsgSerialNo,
				ReplyMsgID:       message.Header.MsgID,
				Result:           0,
			},
		}
		respdata, err = msg8001.Encode()
	}
	g_seq += 1

	if respdata != nil {
		//response
		log.Debug("transerver write response:", hex.EncodeToString(respdata))
		_, err = conn.Write(respdata)
		if err != nil {
			log.Error("write failed ", err)
		}
	}

	return 0
}

func client(ch chan []byte, transAddr string, ctrl_ch chan int, deviceConn *net.TCPConn) {
	g_gwclient_routine++
	for {
		gwConn, err := net.Dial("tcp", transAddr)
		if err != nil {
			log.Warn("init failed", err)
			continue
		}
		tcpGwConn := gwConn.(*net.TCPConn)
		if err != nil {
			log.Info("connect ", transAddr, " fail", err)
		} else {
			log.Info("connect ", transAddr, " ok")
			defer gwConn.Close()
			doTask(tcpGwConn, ch, deviceConn)
		}

		select {
		case _ = <-ctrl_ch:
			log.Info("client ready to exit")
			g_gwclient_routine--
			return
		default:
			log.Debug("nothing todo")
			//nothing todo
		}
		time.Sleep(3 * time.Second)
	}
	g_gwclient_routine--
}

func doTask(gwConn *net.TCPConn, ch chan []byte, deviceConn *net.TCPConn) {
	go receiver(gwConn, deviceConn)
	for {
		//主协程退出不关闭ch，会一直阻塞
		data, ok := <-ch
		if !ok {
			log.Warn("channel closed")
			break
		}

		//转义
		transedData := transform(data)
		log.Debug("begin transfer transformed ", hex.EncodeToString(transedData))
		_, err := gwConn.Write(transedData)
		if err != nil {
			log.Warn("write failed ", hex.EncodeToString(data), err)
			//reconnect need auth, construct one
			message := new(jt808.Message)
			message.Decode(data)
			var imei uint64 = message.Header.Imei
			authMessage := jt808.Message{
				Header: jt808.Header{
					Imei:        imei,
					MsgSerialNo: 1,
				},
				Body: &jt808.T808_0x0102{
					AuthKey: strconv.FormatUint(imei, 10),
				},
			}
			authdata, _ := authMessage.Encode()
			log.Info("build reconnect auth msg:", hex.EncodeToString(authdata))
			//check ch is closed or not
			dataCheck, ok := <-ch
			if !ok {
				log.Warn("channel closed")
				break
			}

			ch <- authdata
			ch <- data

			if dataCheck != nil {
				ch <- dataCheck
			}
			break //break for to reconnect gw
		}
	}
}

type void struct{}

func receiver(gwConn *net.TCPConn, deviceConn *net.TCPConn) {
	g_gwreceiver_routine++
	for {
		err := gwConn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			log.Error("SetReadDeadline failed:", err)
			g_gwreceiver_routine--
			return
		}

		recvBuf := make([]byte, 1024)

		n, err := gwConn.Read(recvBuf[:]) // recv data
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				//log.Println("read timeout:", err)
				continue
			} else {
				log.Warn("read error:", err)
				break
			}
		}
		data := recvBuf[:n]
		//log.Println("response from gw is ", hex.EncodeToString(data))
		message := new(jt808.Message)
		message.Decode(data)
		var member void
		ignore := make(map[int]void)
		ignore[0x8001] = member
		ignore[0x8109] = member
		ignore[0x8100] = member
		_, ok := ignore[int(message.Header.MsgID)]
		if !ok {
			log.Info("write to device from gw ", hex.EncodeToString(data))
			_, err := deviceConn.Write(data)
			if err != nil {
				log.Warn("write to device failed ", err)
				break
			}
		}
	}
	g_gwreceiver_routine--
}
