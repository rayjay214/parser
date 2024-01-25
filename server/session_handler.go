package server

import (
    "github.com/funny/link"
    "github.com/rayjay214/parser/jt808"
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
    }).Debug("[JT/T 808] new session created")

    // 创建Session
    session := newSession(handler.server, sess)
    handler.server.mutex.Lock()
    handler.server.sessions[sess.ID()] = session
    handler.server.mutex.Unlock()
    handler.server.timer.Update(strconv.FormatUint(session.ID(), 10))
    sess.AddCloseCallback(nil, nil, func() {
        handler.server.handleClose(session)
    })

    for {
        // 接收消息
        msg, err := sess.Receive()
        if err != nil {
            sess.Close()
            break
        }

        // 分发消息
        message := msg.(jt808.Message)
        if message.Body == nil || reflect.ValueOf(message.Body).IsNil() {
            session.Reply(&message, jt808.T808_0x8001ResultUnsupported)
            continue
        }

        if !handler.autoMergePacket || !message.Header.Property.IsEnablePacket() {
            session.message(&message)
            handler.server.dispatchMessage(session, &message)
            continue
        }
    }
}
