package server

import (
    "crypto/rsa"
    "github.com/rayjay214/link"
    "github.com/rayjay214/parser/common"
    "github.com/rayjay214/parser/jt808"
    log "github.com/sirupsen/logrus"
    "runtime/debug"
    "strconv"
    "sync"
    "sync/atomic"
    "time"
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
    Protocol uint32

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

// 获取RSA公钥
func (session *Session) GetPublicKey() *rsa.PublicKey {
    codec, ok := session.session.Codec().(*ProtocolCodec)
    if !ok || codec == nil {
        return nil
    }
    return codec.GetPublicKey()
}

// 设置RSA公钥
func (session *Session) SetPublicKey(publicKey *rsa.PublicKey) {
    codec, ok := session.session.Codec().(*ProtocolCodec)
    if !ok || codec == nil {
        return
    }
    codec.SetPublicKey(publicKey)
}

// 发送消息
func (session *Session) Send(entity jt808.Entity) (uint16, error) {
    message := jt808.Message{
        Body: entity,
        Header: jt808.Header{
            MsgID:       entity.MsgID(),
            Imei:        atomic.LoadUint64(&session.imei),
            MsgSerialNo: session.nextID(),
        },
    }

    data, _ := message.Encode()
    log.Printf("send cmd %x", common.GetHex(data))

    err := session.session.Send(message)
    if err != nil {
        return 0, err
    }
    return message.Header.MsgSerialNo, nil
}

// 回复消息
func (session *Session) Reply(msg *jt808.Message, result jt808.Result) (uint16, error) {
    entity := jt808.T808_0x8001{
        ReplyMsgSerialNo: msg.Header.MsgSerialNo,
        ReplyMsgID:       msg.Header.MsgID,
        Result:           result,
    }
    return session.Send(&entity)
}

// 回复注册
func (session *Session) ReplyRegister(msg *jt808.Message) (uint16, error) {
    entity := jt808.T808_0x8100{
        MsgSerialNo: msg.Header.MsgSerialNo,
        AuthKey:     strconv.FormatUint(msg.Header.Imei, 10),
        Result:      0,
    }
    return session.Send(&entity)
}

// 回复短录音
func (session *Session) ReplyShortRecord(pkgNo byte) (uint16, error) {
    entity := jt808.T808_0x8117{
        PkgNo:     pkgNo,
        SessionId: "123454678",
    }
    return session.Send(&entity)
}

// 回复短录音
func (session *Session) ReplyVorRecord(body *jt808.T808_0x0118) (uint16, error) {
    entity := jt808.T808_0x8118{
        PkgNo:     body.PkgNo,
        SessionId: "123454678",
        Time:      body.Time,
    }
    return session.Send(&entity)
}

// 回复校时
func (session *Session) ReplyTime() (uint16, error) {
    now := time.Now()
    entity := jt808.T808_0x8109{
        Year:   uint16(now.Year()),
        Month:  byte(now.Month()),
        Day:    byte(now.Day()),
        Hour:   byte(now.Hour()),
        Minute: byte(now.Minute()),
        Second: byte(now.Second()),
        Result: 0,
    }
    return session.Send(&entity)
}

func (session *Session) Reply8108() (uint16, error) {
    entity := jt808.T808_0x8108{}
    return session.Send(&entity)
}

func (session *Session) Reply8125() (uint16, error) {
    entity := jt808.T808_0x8125{}
    return session.Send(&entity)
}

func (session *Session) Reply8115(sessionId string) (uint16, error) {
    entity := jt808.T808_0x8115{
        SessionId: sessionId,
    }
    return session.Send(&entity)
}

// 发送文本指令
func (session *Session) SendCmd(content string) (uint16, error) {
    entity := jt808.T808_0x8300{
        Flag: 0,
        Text: content,
    }
    return session.Send(&entity)
}

// 开启短录音
func (session *Session) OpenShortRecord(seconds uint64) (uint16, error) {
    entity := jt808.T808_0x8116{
        RecordTime: byte(seconds),
        SessionId:  "12345678",
    }
    return session.Send(&entity)
}

// 声控录音
func (session *Session) VorRecordSwitch(switchs int32) (uint16, error) {
    var entity jt808.T808_0x8103
    if session.Protocol == 1 {
        entity = jt808.T808_0x8103{
            Params: []jt808.Param{
                new(jt808.Param).SetByte(0x0061, byte(switchs)),
            },
        }
    } else {
        entity = jt808.T808_0x8103{
            Params: []jt808.Param{
                new(jt808.Param).SetByte(0xf114, byte(switchs)),
            },
        }
    }

    return session.Send(&entity)
}

// 发起请求
func (session *Session) Request(entity jt808.Entity, cb func(answer *jt808.Message)) (uint16, error) {
    serialNo, err := session.Send(entity)
    if err != nil {
        return 0, err
    }

    if cb != nil {
        session.addRequestContext(requestContext{
            msgID:    uint16(entity.MsgID()),
            serialNo: serialNo,
            callback: cb,
        })
    }
    return serialNo, nil
}

// 关闭连接
func (session *Session) Close() error {
    return session.session.Close()
}

// 获取消息ID
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

// 消息接收事件
func (session *Session) message(message *jt808.Message) {
    if message.Header.Imei > 0 {
        old := atomic.LoadUint64(&session.imei)
        if old != 0 && old != message.Header.Imei {
            log.WithFields(log.Fields{
                "id":  session.ID(),
                "old": old,
                "new": message.Header.Imei,
            }).Warn("[JT/T 808] terminal Imei is inconsistent")
        }
        atomic.StoreUint64(&session.imei, message.Header.Imei)
    }

    var msgSerialNo uint16
    switch message.Header.MsgID {
    case jt808.MsgT808_0x0001:
        // 终端通用应答
        msgSerialNo = message.Body.(*jt808.T808_0x0001).ReplyMsgSerialNo
    case jt808.MsgT808_0x0104:
        // 查询终端参数应答
        msgSerialNo = message.Body.(*jt808.T808_0x0104).ReplyMsgSerialNo
    case jt808.MsgT808_0x0201:
        // 位置信息查询应答
        msgSerialNo = message.Body.(*jt808.T808_0x0201).ReplyMsgSerialNo
    case jt808.MsgT808_0x0302:
        // 提问应答
        msgSerialNo = message.Body.(*jt808.T808_0x0302).ReplyMsgSerialNo
    case jt808.MsgT808_0x0500:
        // 车辆控制应答
        msgSerialNo = message.Body.(*jt808.T808_0x0500).ReplyMsgSerialNo
    case jt808.MsgT808_0x0700:
        // 行驶记录数据上传
        msgSerialNo = message.Body.(*jt808.T808_0x0700).ReplyMsgSerialNo
    case jt808.MsgT808_0x0802:
        // 存储多媒体数据检索应答
        msgSerialNo = message.Body.(*jt808.T808_0x0802).ReplyMsgSerialNo
    case jt808.MsgT808_0x0805:
        // 摄像头立即拍摄命令应答
        msgSerialNo = message.Body.(*jt808.T808_0x0805).ReplyMsgSerialNo
    }
    if msgSerialNo == 0 {
        return
    }

    ctx, ok := session.takeRequestContext(msgSerialNo)
    if ok {
        defer func() {
            if err := recover(); err != nil {
                debug.PrintStack()
            }
        }()
        ctx.callback(message)
    }
}

// 添加请求上下文
func (session *Session) addRequestContext(ctx requestContext) {
    session.mux.Lock()
    defer session.mux.Unlock()

    for idx, item := range session.requests {
        if item.msgID == ctx.msgID {
            session.requests[idx] = ctx
            return
        }
    }
    session.requests = append(session.requests, ctx)
}

// 取出请求上下文
func (session *Session) takeRequestContext(msgSerialNo uint16) (requestContext, bool) {
    session.mux.Lock()
    defer session.mux.Unlock()

    for idx, item := range session.requests {
        if item.serialNo == msgSerialNo {
            session.requests[idx] = session.requests[len(session.requests)-1]
            session.requests = session.requests[:len(session.requests)-1]
            return item, true
        }
    }
    return requestContext{}, false
}
