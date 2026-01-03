package hl3g_base

import (
	"github.com/rayjay214/link"
	"github.com/rayjay214/parser/protocol/hl3g"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
)

// Session处理
type sessionHandler struct {
	server          *Server
	autoMergePacket bool
}

func (handler sessionHandler) HandleSession(sess *link.Session) {
	log.WithFields(log.Fields{
		"id": sess.ID(),
	}).Info("[hl3g] new session created")

	var session *Session

	for {
		// 接收消息
		msg, err := sess.Receive()
		if err != nil {
			sess.Close()
			break
		}

		// 分发消息
		message := msg.(hl3g.Message)
		if message.Body == nil || reflect.ValueOf(message.Body).IsNil() {
			if session != nil {
				//session.Reply(&message, jt808.T808_0x8001ResultUnsupported)
			}
			continue
		}

		if message.Header.Proto == "LK2" {
			imei, _ := strconv.ParseUint(message.Header.Imei, 10, 64)
			deviceInfo, _ := storage.GetDevice(imei)
			if len(deviceInfo) == 0 {
				log.Warnf("imei %v not exist", message.Header.Imei)
				sess.Close()
				break
			}

			session = newSession(handler.server, sess)
			handler.server.mutex.Lock()
			delete(handler.server.sessions, imei)
			handler.server.sessions[imei] = session
			session.imei = imei
			session.UserData = make(map[string]interface{}, 8)
			session.Protocol, _ = strconv.Atoi(deviceInfo["protocol"])
			handler.server.mutex.Unlock()
			handler.server.timer.Update(strconv.FormatUint(session.ID(), 10))
			sess.AddCloseCallback(nil, nil, func() {
				handler.server.handleClose(session)
			})
		}

		handler.server.dispatchMessage(session, &message)
	}
}
