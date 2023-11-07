package ipc

// 消息ID枚举
type MsgID uint16

const PrefixID = byte(0x86)

const (
    // 终端--》平台
    // 鉴权
    Msg_0x0001 MsgID = 0x0001
    // 心跳
    Msg_0x0002 MsgID = 0x0002
    // 上报公网信息
    Msg_0x0003 MsgID = 0x0003
    // 应答文本消息
    Msg_0x1300 MsgID = 0x1300

    // 平台--》终端
    // 通用回复
    Msg_0x8001 MsgID = 0x8001
    // 文本消息下发
    Msg_0x8300 MsgID = 0x8300
)

// 消息实体映射
var entityMapper = map[MsgID]func() Entity{
    Msg_0x0001: func() Entity {
        return new(Body_0x0001)
    },
    Msg_0x0002: func() Entity {
        return new(Body_0x0002)
    },
    Msg_0x0003: func() Entity {
        return new(Body_0x0003)
    },
    /*
    	Msg_0x1300: func() Entity {
    		return new(Body_0x1300)
    	},
    	Msg_0x8001: func() Entity {
    		return new(Body_0x8001)
    	},
    	Msg_0x8300: func() Entity {
    		return new(Body_0x8300)
    	},
    */
}
