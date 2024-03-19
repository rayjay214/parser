package server

import (
	"github.com/rayjay214/link"
	"github.com/rayjay214/parser/jt808"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
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

	for {
		// 接收消息
		msg, err := sess.Receive()
		if err != nil {
			log.Warnf("%v receive err %v", sess.ID(), err)
			//实测发现,timeout之后,设备并不会重连,也就是不会发0102,所以session无法创建，而原先的session如果被删除之后,也就无法下发指令
			if strings.Contains(err.Error(), "timed out") {
				log.Warnf("%v timeout， wait", session.ID())
				continue
			}
			sess.Close()
			break
			//continue //这里不关闭链接，依赖定时器去清理链接
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

		if message.Header.MsgID == jt808.MsgT808_0x0002 { //有时间服务器连接断了，删除session了，但是设备不知道，没有重连，没有上报0102
			_, ok := handler.server.sessions[message.Header.Imei]
			if !ok {
				log.Warnf("%v session lost, new one", message.Header.Imei)
				deviceInfo, _ := storage.GetDevice(message.Header.Imei)
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
		}

		/*
			if message.Header.MsgID != jt808.MsgT808_0x0117 && message.Header.MsgID != jt808.MsgT808_0x0118 {
				d, _ := message.Encode()
				log.WithFields(log.Fields{
					"imei": message.Header.Imei,
				}).Infof("receive msg %x", common.GetHex(d))
			}
		*/

		//session.message(&message)
		handler.server.dispatchMessage(session, &message)
	}
}
