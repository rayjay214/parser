package gmconf

// 消息ID枚举
type MsgID uint8

const PrefixID = byte(0x66)

const (
	Msg_0x13 MsgID = 0x13
	Msg_0x14 MsgID = 0x14
	Msg_0x92 MsgID = 0x92
)

// 消息实体映射
var entityMapper = map[uint8]func() Entity{
	uint8(Msg_0x13): func() Entity {
		return new(Body_0x13)
	},
	uint8(Msg_0x14): func() Entity {
		return new(Body_0x14)
	},
	uint8(Msg_0x92): func() Entity {
		return new(Body_0x92)
	},
}
