package tq

// 消息ID枚举
type MsgID string

const (
    // 心跳
    Msg_V1 MsgID = "V1"
    // 地址请求
    Msg_V2 MsgID = "V2"
    // CMD返回
    Msg_V4 MsgID = "V4"
    // 带里程
    Msg_V5 MsgID = "V5"
    // 带ICCID
    Msg_V6 MsgID = "V6"
    // 指令回复
    Msg_I1 MsgID = "I1"
    // 通用回复
    Msg_R12 MsgID = "R12"
    // 正常数据包
    Msg_Normal MsgID = "Normal"
)

// 消息实体映射
var entityMapper = map[string]func() Entity{
    string(Msg_V1): func() Entity {
        return new(TQ_V1)
    },
    string(Msg_V2): func() Entity {
        return new(TQ_V2)
    },
    string(Msg_V4): func() Entity {
        return new(TQ_V4)
    },
    string(Msg_V5): func() Entity {
        return new(TQ_V5)
    },
    string(Msg_V6): func() Entity {
        return new(TQ_V6)
    },
    string(Msg_I1): func() Entity {
        return new(TQ_I1)
    },
    string(Msg_R12): func() Entity {
        return new(TQ_R12)
    },
    string(Msg_Normal): func() Entity {
        return new(TQ_Normal)
    },
}
