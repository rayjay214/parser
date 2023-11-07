package th

// 消息ID枚举
type MsgID string

const (
    Msg_GPS MsgID = "Gpslocation"

    Msg_LBS MsgID = "Lbslocation"

    Msg_ONLINE MsgID = "#006"

    Msg_FOTA MsgID = "#007"
)

// 消息实体映射
var entityMapper = map[string]func() Entity{
    string(Msg_GPS): func() Entity {
        return new(TH_GPS)
    },
    string(Msg_LBS): func() Entity {
        return new(TH_LBS)
    },
    string(Msg_ONLINE): func() Entity {
        return new(TH_ONLINE)
    },
    string(Msg_FOTA): func() Entity {
        return new(TH_FOTA)
    },
}
