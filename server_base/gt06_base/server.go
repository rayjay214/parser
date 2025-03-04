package gt06_base

import (
	"errors"
	"fmt"
	"github.com/rayjay214/link"
	"github.com/rayjay214/parser/protocol/gt06"
	"github.com/rayjay214/parser/server_base/common"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"net"
	"runtime/debug"
	"strconv"
	"sync"
)

// 服务器选项
type Options struct {
	Keepalive       int64
	AutoMergePacket bool
	CloseHandler    func(*Session)
}

// 协议服务器
type Server struct {
	server  *link.Server
	handler sessionHandler
	timer   *common.CountdownTimer

	mutex    sync.Mutex
	sessions map[uint64]*Session

	closeHandler    func(*Session)
	messageHandlers sync.Map
}

type MessageHandler func(*Session, *gt06.Message)

// 创建服务
func NewServer(options Options) (*Server, error) {
	if options.Keepalive <= 0 {
		options.Keepalive = 60
	}

	server := Server{
		closeHandler: options.CloseHandler,
		sessions:     make(map[uint64]*Session),
	}
	server.handler.server = &server
	server.handler.autoMergePacket = options.AutoMergePacket
	server.timer = common.NewCountdownTimer(options.Keepalive, server.handleReadTimeout)
	return &server, nil
}

// 运行服务
func (server *Server) Run(network string, port int) error {
	if server.server != nil {
		return errors.New("server already running")
	}

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listen, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	p := Protocol{}
	server.server = link.NewServer(listen, &p, 24, server.handler)
	log.Infof("[JT/T 808] protocol server started on %s", address)
	return server.server.Serve()
}

// 停止服务
func (server *Server) Stop() {
	if server.server != nil {
		server.server.Stop()
		server.server = nil
	}
}

// 获取Session
func (server *Server) GetSession(id uint64) (*Session, bool) {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	session, ok := server.sessions[id]
	if !ok {
		return nil, false
	}
	return session, true
}

// 添加消息处理
func (server *Server) AddHandler(msgID gt06.MsgID, handler MessageHandler) {
	if handler != nil {
		server.messageHandlers.Store(msgID, handler)
	}
}

// 处理关闭
func (server *Server) handleClose(session *Session) {
	server.mutex.Lock()
	delete(server.sessions, session.ID())
	server.mutex.Unlock()

	server.timer.Remove(strconv.FormatUint(session.ID(), 10))
	if server.closeHandler != nil {
		func() {
			defer func() {
				if err := recover(); err != nil {
					debug.PrintStack()
				}
			}()
			server.closeHandler(session)
		}()
	}

	info := map[string]interface{}{
		"state": "1",
	}
	storage.SetRunInfo(session.ID(), info)
	storage.InsertOfflineAlarm(session.ID())

	log.WithFields(log.Fields{
		"id": session.ID(),
	}).Info("[gt06] session closed")
}

// 处理读超时
func (server *Server) handleReadTimeout(key string) {
	sessionID, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		return
	}

	session, ok := server.GetSession(sessionID)
	if !ok {
		return
	}
	session.Close()

	info := map[string]interface{}{
		"state": "1",
	}
	storage.SetRunInfo(sessionID, info)
	storage.InsertOfflineAlarm(sessionID)

	log.WithFields(log.Fields{
		"id": sessionID,
	}).Info("[gt06] session read timeout")
}

// 分派消息
func (server *Server) dispatchMessage(session *Session, message *gt06.Message) {
	id := message.Body.MsgID()
	fmt.Println(id)
	handler, ok := server.messageHandlers.Load(id)
	if !ok {
		log.WithFields(log.Fields{
			"id": fmt.Sprintf("0x%x", message.Body.MsgID),
		}).Info("[gt06] dispatch message canceled, handler not found")
		return
	}

	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	server.timer.Update(strconv.FormatUint(session.ID(), 10))
	handler.(MessageHandler)(session, message)
}
