package ota

// 消息ID枚举
type MsgID uint8

const PrefixID = byte(0x68)

const (
    // 终端向服务器查询是否有升级
    Msg_0x11 MsgID = 0x11
    // 服务器向终端返回是否有升级
    Msg_0x91 MsgID = 0x91
    // 终端向服务器上报升级结果
    Msg_0x12 MsgID = 0x12
)

// 消息实体映射
var entityMapper = map[uint8]func() Entity{
    uint8(Msg_0x11): func() Entity {
        return new(Body_0x11)
    },
    uint8(Msg_0x91): func() Entity {
        return new(Body_0x91)
    },
    uint8(Msg_0x12): func() Entity {
        return new(Body_0x12)
    },
}
