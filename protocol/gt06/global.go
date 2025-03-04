package gt06

// 消息ID枚举
type MsgID uint8

const (
	// 登录
	Msg_0x01 MsgID = 0x01
	// 心跳
	Msg_0x13 MsgID = 0x13
	// 定位包
	Msg_0x22 MsgID = 0x22
	// 告警包
	Msg_0x26 MsgID = 0x26
	// 通用信息传输
	Msg_0x94 MsgID = 0x94
	// 在线指令发送
	Msg_0x80 MsgID = 0x80
	// 在线指令回复
	Msg_0x21 MsgID = 0x21
	// LBS地址请求包
	Msg_0x17 MsgID = 0x17
	// 地址请求中文回复
	Msg_0x17_AddrResp MsgID = 0xF0
	// 告警请求中文回复
	Msg_0x17_AlarmResp MsgID = 0xF1
	// GPS地址请求包
	Msg_0x2a MsgID = 0x2a
	// gt06e定位包
	Msg_0x12 MsgID = 0x12
	// 告警包
	Msg_0x16 MsgID = 0x16
	// 多基站包
	Msg_0x28 MsgID = 0x28
	// 多基站包(gt06e)
	Msg_0x18 MsgID = 0x18
	// 在线指令回复(gt06e)
	Msg_0x15 MsgID = 0x15
	// 校时包
	Msg_0x8a MsgID = 0x8a
	// 4G定位包
	Msg_0xa0 MsgID = 0xa0
	// 4G报警包
	Msg_0xa3 MsgID = 0xa3
	// 4GLBS多基站包
	Msg_0xa1 MsgID = 0xa1
)

// 消息实体映射
var entityMapper = map[uint8]func() Entity{
	uint8(Msg_0x01): func() Entity {
		return new(Kks_0x01)
	},
	uint8(Msg_0x13): func() Entity {
		return new(Kks_0x13)
	},
	uint8(Msg_0x94): func() Entity {
		return new(Kks_0x94)
	},
	uint8(Msg_0x80): func() Entity {
		return new(Kks_0x80)
	},
	uint8(Msg_0x21): func() Entity {
		return new(Kks_0x21)
	},
	uint8(Msg_0x12): func() Entity {
		return new(Kks_0x12)
	},
	uint8(Msg_0x16): func() Entity {
		return new(Kks_0x16)
	},
	uint8(Msg_0x28): func() Entity {
		return new(Kks_0x28)
	},
	uint8(Msg_0x15): func() Entity {
		return new(Kks_0x15)
	},
	uint8(Msg_0xa1): func() Entity {
		return new(Kks_0xa1)
	},
}
