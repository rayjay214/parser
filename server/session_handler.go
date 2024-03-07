package server

import (
	"github.com/rayjay214/link"
	"github.com/rayjay214/parser/common"
	"github.com/rayjay214/parser/jt808"
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
	}).Info("[JT/T 808] new session created")

	var session *Session

	// 创建Session
	/*
	   session := newSession(handler.server, sess)
	   handler.server.mutex.Lock()
	   handler.server.sessions[sess.ID()] = session
	   handler.server.mutex.Unlock()
	   handler.server.timer.Update(strconv.FormatUint(session.ID(), 10))
	   sess.AddCloseCallback(nil, nil, func() {
	       handler.server.handleClose(session)
	   })
	*/

	for {
		// 接收消息
		msg, err := sess.Receive()
		if err != nil {
			log.Warnf("receive err %v", err)
			sess.Close()
			break
		}

		// 分发消息
		message := msg.(jt808.Message)
		if message.Body == nil || reflect.ValueOf(message.Body).IsNil() {
			if session != nil {
				session.Reply(&message, jt808.T808_0x8001ResultUnsupported)
			}
			continue
		}

		//_, ok := handler.server.sessions[message.Header.Imei]
		if message.Header.MsgID == jt808.MsgT808_0x0100 || message.Header.MsgID == jt808.MsgT808_0x0102 {
			deviceInfo, _ := storage.GetDevice(message.Header.Imei)
			if len(deviceInfo) == 0 {
				log.Warnf("imei %v not exist", message.Header.Imei)
				sess.Close()
				break
			}

			session = newSession(handler.server, sess)
			handler.server.mutex.Lock()
			delete(handler.server.sessions, message.Header.Imei)
			handler.server.sessions[message.Header.Imei] = session
			session.imei = message.Header.Imei
			session.UserData = make(map[string]interface{}, 8)
			session.Protocol, _ = strconv.Atoi(deviceInfo["protocol"])
			handler.server.mutex.Unlock()
			handler.server.timer.Update(strconv.FormatUint(session.ID(), 10))
			sess.AddCloseCallback(nil, nil, func() {
				handler.server.handleClose(session)
			})
		}

		if message.Header.MsgID != jt808.MsgT808_0x0117 && message.Header.MsgID != jt808.MsgT808_0x0118 {
			d, _ := message.Encode()
			log.WithFields(log.Fields{
				"imei": message.Header.Imei,
			}).Infof("receive msg %x", common.GetHex(d))
		}

		//session.message(&message)
		handler.server.dispatchMessage(session, &message)
	}
}
