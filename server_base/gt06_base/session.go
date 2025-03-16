package gt06_base

import (
	"github.com/rayjay214/link"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/rayjay214/parser/protocol/gt06"
	"github.com/rayjay214/parser/protocol/jt808"
	log "github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
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

// 发送消息
func (session *Session) Send(entity gt06.Entity) (uint16, error) {

	message := gt06.Message{
		Body: entity,
		Header: gt06.Header{
			Prefix: 0x7878,
		},
	}

	data, _ := message.Encode()
	log.Infof("rayjay07 send msg %x", common.GetHex(data))
	err := session.session.Send(message)
	if err != nil {
		return 0, err
	}
	return entity.GetSeqNo(), nil
}

func (session *Session) SendCmd(content string) (uint16, error) {
	seqNo := session.nextID()
	entity := gt06.Kks_0x80{
		Proto:   0x80,
		Content: content,
		SeqNo:   seqNo,
		SysFlag: uint32(seqNo),
	}
	return session.Send(&entity)
}

func (session *Session) CommonReply(proto byte) (uint16, error) {
	entity := gt06.KksResponse{
		Proto: proto,
		SeqNo: session.nextID(),
	}
	return session.Send(&entity)
}

func (session *Session) Close() error {
	return session.session.Close()
}

func (session *Session) nextID() uint16 {
	var id uint32
	for {
		id = atomic.LoadUint32(&session.next)
		if id == 0xff {
			if atomic.CompareAndSwapUint32(&session.next, id, 1) {
				id = 1
				break
			}
		} else if atomic.CompareAndSwapUint32(&session.next, id, id+1) {
			id += 1
			break
		}
	}
	return uint16(id)
}
