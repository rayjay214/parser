package hl3g_base

import (
	"github.com/rayjay214/link"
	"github.com/rayjay214/parser/protocol/hl3g"
	"github.com/rayjay214/parser/protocol/jt808"
	log "github.com/sirupsen/logrus"
	"sync"
)

// 请求上下文
type requestContext struct {
	msgID    uint16
	serialNo uint16
	callback func(answer *jt808.Message)
}

// 终端会话
type Session struct {
	next    uint32
	imei    uint64
	server  *Server
	session *link.Session

	mux      sync.Mutex
	requests []requestContext
	Protocol int

	UserData map[string]interface{}
}

// 创建Session
func newSession(server *Server, sess *link.Session) *Session {
	return &Session{
		server:  server,
		session: sess,
	}
}

// 获取ID
func (session *Session) ID() uint64 {
	return session.imei
}

// 获取服务实例
func (session *Session) GetServer() *Server {
	return session.server
}

func (session *Session) CommonReply(imei, proto string) (uint16, error) {
	var message hl3g.Message
	var header hl3g.Header
	header.Prefix = "[3G"
	header.Imei = imei
	header.MsgLen = "0000"
	header.Proto = proto

	message.Header = header

	msg, _ := message.Encode()
	log.Infof("%v send msg %s", imei, string(msg))

	err := session.session.Send(message)
	if err != nil {
		return 0, err
	}
	return 0, err
}

func (session *Session) UploadReply(imei, proto string, interval string) (uint16, error) {
	var message hl3g.Message
	var header hl3g.Header
	header.Prefix = "[3G"
	header.Imei = imei
	header.MsgLen = "0000"
	header.Proto = proto

	message.Header = header

	var body hl3g.HL3G_UPLOAD
	body.Interval = interval

	message.Body = &body

	msg, _ := message.Encode()
	log.Infof("%v send msg %s", imei, string(msg))

	err := session.session.Send(message)
	if err != nil {
		return 0, err
	}
	return 0, err
}

func (session *Session) ModeWorkReply(imei, proto string, mode string) (uint16, error) {
	var message hl3g.Message
	var header hl3g.Header
	header.Prefix = "[3G"
	header.Imei = imei
	header.MsgLen = "0000"
	header.Proto = proto

	message.Header = header

	var body hl3g.HL3G_MODEWORK
	body.Mode = mode

	message.Body = &body

	msg, _ := message.Encode()
	log.Infof("%v send msg %s", imei, string(msg))

	err := session.session.Send(message)
	if err != nil {
		return 0, err
	}
	return 0, err
}

func (session *Session) CloseModeReply(imei, proto string, content string) (uint16, error) {
	var message hl3g.Message
	var header hl3g.Header
	header.Prefix = "[3G"
	header.Imei = imei
	header.MsgLen = "0000"
	header.Proto = proto

	message.Header = header

	var body hl3g.HL3G_CLOSEMODE
	body.Content = content

	message.Body = &body

	msg, _ := message.Encode()
	log.Infof("%v send msg %s", imei, string(msg))

	err := session.session.Send(message)
	if err != nil {
		return 0, err
	}
	return 0, err
}

func (session *Session) TcReply(imei, proto string, content string) (uint16, error) {
	var message hl3g.Message
	var header hl3g.Header
	header.Prefix = "[3G"
	header.Imei = imei
	header.MsgLen = "0000"
	header.Proto = proto

	message.Header = header

	var body hl3g.HL3G_TC
	body.Content = content

	message.Body = &body

	msg, _ := message.Encode()
	log.Infof("%v send msg %s", imei, string(msg))

	err := session.session.Send(message)
	if err != nil {
		return 0, err
	}
	return 0, err
}

func (session *Session) Gs1Reply(imei, proto string, lat, lng, now string) (uint16, error) {
	var message hl3g.Message
	var header hl3g.Header
	header.Prefix = "[3G"
	header.Imei = imei
	header.MsgLen = "0000"
	header.Proto = proto

	message.Header = header

	var body hl3g.HL3G_GS
	body.Lat = lat
	body.Lng = lng
	body.Time = now

	message.Body = &body

	msg, _ := message.Encode()
	log.Infof("%v send msg %s", imei, string(msg))

	err := session.session.Send(message)
	if err != nil {
		return 0, err
	}
	return 0, err
}

func (session *Session) Close() error {
	return session.session.Close()
}
