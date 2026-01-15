package hl3g

// 消息ID枚举
type MsgID string

const (
	Msg_LK2       MsgID = "LK2"
	Msg_GS1       MsgID = "GS1"
	Msg_GS        MsgID = "GS"
	Msg_CCID      MsgID = "CCID"
	Msg_UD        MsgID = "UD"
	Msg_UD2       MsgID = "UD2"
	Msg_AL        MsgID = "AL"
	Msg_UPLOAD    MsgID = "UPLOAD"
	Msg_IP        MsgID = "IP"
	Msg_FACTORY   MsgID = "FACTORY"
	Msg_VERNO     MsgID = "VERNO"
	Msg_RESET     MsgID = "RESET"
	Msg_STU       MsgID = "STU"
	Msg_MODEWORK  MsgID = "MODEWORK"
	Msg_TC        MsgID = "TC"
	Msg_CLOSEMODE MsgID = "CLOSEMODE"
)

// 消息实体映射
var entityMapper = map[string]func() Entity{
	string(Msg_LK2): func() Entity {
		return new(HL3G_LK2)
	},
	string(Msg_GS1): func() Entity {
		return new(HL3G_GS1)
	},
	string(Msg_GS): func() Entity {
		return new(HL3G_GS)
	},
	string(Msg_CCID): func() Entity {
		return new(HL3G_CCID)
	},
	string(Msg_UD): func() Entity {
		return new(HL3G_UD)
	},
	string(Msg_UD2): func() Entity {
		return new(HL3G_UD2)
	},
	string(Msg_AL): func() Entity {
		return new(HL3G_AL)
	},
	string(Msg_IP): func() Entity {
		return new(HL3G_LK2)
	},
	string(Msg_FACTORY): func() Entity {
		return new(HL3G_FACTORY)
	},
	string(Msg_VERNO): func() Entity {
		return new(HL3G_VERNO)
	},
	string(Msg_RESET): func() Entity {
		return new(HL3G_RESET)
	},
	string(Msg_UPLOAD): func() Entity {
		return new(HL3G_UPLOAD)
	},
	string(Msg_STU): func() Entity {
		return new(HL3G_STU)
	},
	string(Msg_MODEWORK): func() Entity {
		return new(HL3G_MODEWORK)
	},
	string(Msg_TC): func() Entity {
		return new(HL3G_TC)
	},
	string(Msg_CLOSEMODE): func() Entity {
		return new(HL3G_CLOSEMODE)
	},
}
