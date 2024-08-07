package extra

// 附加消息类型
type Type byte

const (
	// 里程
	TypeExtra_0x01 Type = 0x01
	// 油量
	TypeExtra_0x02 Type = 0x02
	// 状态，电量
	TypeExtra_0x04 Type = 0x04
	// 1:设备运动 0：设备静止
	TypeExtra_0x05 Type = 0x05
	// 1:该地位要发短信使用  0： 表示该定位不用发短信
	TypeExtra_0x06 Type = 0x06
	// 无线通信网络信号强度
	TypeExtra_0x30 Type = 0x30
	// GNSS定位卫星数
	TypeExtra_0x31 Type = 0x31
	// 基站信息上传
	TypeExtra_0xe1 Type = 0xe1
	// 状态，电量
	TypeExtra_0xe4 Type = 0xe4
	// 1:设备运动 0：设备静止
	TypeExtra_0xe5 Type = 0xe5
	// 1:该地位要发短信使用  0： 表示该定位不用发短信
	TypeExtra_0xe6 Type = 0xe6
	// 状态扩展
	TypeExtra_0xe7 Type = 0xe7
	// 基站信息
	TypeExtra_0xeb Type = 0xeb
	// wifi信息
	TypeExtra_0xec Type = 0xec
	// cdma信息
	TypeExtra_0xed Type = 0xed
	// 4g基站信息
	TypeExtra_0xee Type = 0xee
	// 外电电压
	TypeExtra_0xf0 Type = 0xf0
	// 当前电流
	TypeExtra_0xf1 Type = 0xf1
	// GPS卫星个数
	TypeExtra_0xf2 Type = 0xf2
	// 北斗卫星个数
	TypeExtra_0xf3 Type = 0xf3
	// 格洛纳斯卫星个数
	TypeExtra_0xf4 Type = 0xf4
	// 信号类型
	TypeExtra_0xf5 Type = 0xf5
	// Imei上传
	TypeExtra_0xf6 Type = 0xf6
	// Imei上传
	TypeExtra_0x51 Type = 0x51
	// Apn上传
	TypeExtra_0xf7 Type = 0xf7
	//fg add
    // 主电源电压
	TypeExtra_0x2b Type = 0x2b
	// 运动或静止
	TypeExtra_0x32 Type = 0x32
	// iccid
	TypeExtra_0xb2 Type = 0xb2
	// 移文基站
	TypeExtra_0x5d Type = 0x5d
	// 移文wifi
	TypeExtra_0x54 Type = 0x54
	// 计步步数
	TypeExtra_0xf8 Type = 0xf8
	// 是否上报假点
	TypeExtra_0xf9 Type = 0xf9
)

// 消息实体映射
var entityMapper = map[byte]func() Entity{
	byte(TypeExtra_0x01): func() Entity {
		return new(Extra_0x01)
	},
	byte(TypeExtra_0x02): func() Entity {
		return new(Extra_0x02)
	},
	byte(TypeExtra_0x04): func() Entity {
		return new(Extra_0x04)
	},
	byte(TypeExtra_0x05): func() Entity {
		return new(Extra_0x05)
	},
	byte(TypeExtra_0x06): func() Entity {
		return new(Extra_0x06)
	},
	byte(TypeExtra_0x30): func() Entity {
		return new(Extra_0x30)
	},
	byte(TypeExtra_0x31): func() Entity {
		return new(Extra_0x31)
	},
	byte(TypeExtra_0xe1): func() Entity {
		return new(Extra_0xe1)
	},
	byte(TypeExtra_0xe4): func() Entity {
		return new(Extra_0xe4)
	},
	byte(TypeExtra_0xe5): func() Entity {
		return new(Extra_0xe5)
	},
	byte(TypeExtra_0xe6): func() Entity {
		return new(Extra_0xe6)
	},
	byte(TypeExtra_0xe7): func() Entity {
		return new(Extra_0xe7)
	},
	byte(TypeExtra_0xeb): func() Entity {
		return new(Extra_0xeb)
	},
	byte(TypeExtra_0xec): func() Entity {
		return new(Extra_0xec)
	},
	byte(TypeExtra_0xed): func() Entity {
		return new(Extra_0xed)
	},
	byte(TypeExtra_0xee): func() Entity {
		return new(Extra_0xee)
	},
	byte(TypeExtra_0xf0): func() Entity {
		return new(Extra_0xf0)
	},
	byte(TypeExtra_0xf1): func() Entity {
		return new(Extra_0xf1)
	},
	byte(TypeExtra_0xf2): func() Entity {
		return new(Extra_0xf2)
	},
	byte(TypeExtra_0xf3): func() Entity {
		return new(Extra_0xf3)
	},
	byte(TypeExtra_0xf4): func() Entity {
		return new(Extra_0xf4)
	},
	byte(TypeExtra_0xf5): func() Entity {
		return new(Extra_0xf5)
	},
	byte(TypeExtra_0xf6): func() Entity {
		return new(Extra_0xf6)
	},
	byte(TypeExtra_0x51): func() Entity {
		return new(Extra_0x51)
	},
	byte(TypeExtra_0xf7): func() Entity {
		return new(Extra_0xf7)
	},
	byte(TypeExtra_0x5d): func() Entity {
		return new(Extra_0x5d)
	},
	byte(TypeExtra_0x54): func() Entity {
		return new(Extra_0x54)
	},
	byte(TypeExtra_0xf8): func() Entity {
		return new(Extra_0xf8)
	},
	byte(TypeExtra_0xf9): func() Entity {
		return new(Extra_0xf9)
	},
    //fg add
    /*
	byte(TypeExtra_0x2b): func() Entity {
		return new(Extra_0x2b)
	},
	byte(TypeExtra_0x5d): func() Entity {
		return new(Extra_0x5d)
	},
	byte(TypeExtra_0x32): func() Entity {
		return new(Extra_0x32)
	},
	byte(TypeExtra_0xb2): func() Entity {
		return new(Extra_0xb2)
	},
    */
}

// 类型注册
func Register(typ byte, creator func() Entity) {
	entityMapper[typ] = creator
}
