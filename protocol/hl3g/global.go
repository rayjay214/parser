package hl3g

// 消息ID枚举
type MsgID string

const (
	Msg_LK2     MsgID = "LK2"
	Msg_GS1     MsgID = "GS1"
	Msg_GS      MsgID = "GS"
	Msg_CCID    MsgID = "CCID"
	Msg_UD      MsgID = "UD"
	Msg_UD2     MsgID = "UD2"
	Msg_AL      MsgID = "AL"
	Msg_UPLOAD  MsgID = "UPLOAD"
	Msg_IP      MsgID = "IP"
	Msg_FACTORY MsgID = "FACTORY"
	Msg_VERNO   MsgID = "VERNO"
	Msg_RESET   MsgID = "RESET"
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
		return new(HL3G_LK2)
	},
	string(Msg_VERNO): func() Entity {
		return new(HL3G_LK2)
	},
	string(Msg_RESET): func() Entity {
		return new(HL3G_LK2)
	},
	string(Msg_UPLOAD): func() Entity {
		return new(HL3G_LK2)
	},
}
