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
	// 通知开启转发
	Msg_0x0004 MsgID = 0x0004
	// 应答文本消息
	Msg_0x1300 MsgID = 0x1300
	// 应答指令
	Msg_0x1301 MsgID = 0x1301
	// 请求时间
	Msg_0x0021 MsgID = 0x0021
	// 返回文件列表
	Msg_0x1201 MsgID = 0x1201
	// 确认返回文件
	Msg_0x1202 MsgID = 0x1202
	// 通知服务器开始上传文件
	Msg_0x1203 MsgID = 0x1203
	// 上传文件
	Msg_0x1204 MsgID = 0x1204
	// 上传文件结果
	Msg_0x1205 MsgID = 0x1205
	// 返回月中有视频的日期
	Msg_0x1206 MsgID = 0x1206
	// 取消上传文件回复
	Msg_0x1207 MsgID = 0x1207

	// 平台--》终端
	// 鉴权回复
	Msg_0x8000 MsgID = 0x8000
	// 通用回复
	Msg_0x8001 MsgID = 0x8001
	// 文本消息下发
	Msg_0x8300 MsgID = 0x8300
	// 上报app的公网信息
	Msg_0x8101 MsgID = 0x8101
	// app心跳
	Msg_0x8002 MsgID = 0x8002
	// app通知p2p连接已建立
	Msg_0x8003 MsgID = 0x8003
	// app通知设备播放事件
	Msg_0x8004 MsgID = 0x8004
	// app通知设备转发
	Msg_0x8005 MsgID = 0x8005
	// 指令下发
	Msg_0x8301 MsgID = 0x8301
	// 校时回复
	Msg_0x8021 MsgID = 0x8021
	// 请求终端录像文件列表
	Msg_0x8201 MsgID = 0x8201
	// 请求某个录像文件
	Msg_0x8202 MsgID = 0x8202
	// 回复1203
	Msg_0x8203 MsgID = 0x8203
	// 确认收到1204
	Msg_0x8204 MsgID = 0x8204
	// 回复1205
	Msg_0x8205 MsgID = 0x8205
	// 请求月中有视频的日期
	Msg_0x8206 MsgID = 0x8206
	// 取消文件上传
	Msg_0x8207 MsgID = 0x8207
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
	Msg_0x0004: func() Entity {
		return new(Body_0x0004)
	},
	Msg_0x1300: func() Entity {
		return new(Body_0x1300)
	},
	Msg_0x1301: func() Entity {
		return new(Body_0x1301)
	},
	Msg_0x0021: func() Entity {
		return new(Body_0x0021)
	},
	Msg_0x1201: func() Entity {
		return new(Body_0x1201)
	},
	Msg_0x1202: func() Entity {
		return new(Body_0x1202)
	},
	Msg_0x1203: func() Entity {
		return new(Body_0x1203)
	},
	Msg_0x1204: func() Entity {
		return new(Body_0x1204)
	},
	Msg_0x1205: func() Entity {
		return new(Body_0x1205)
	},
	Msg_0x1206: func() Entity {
		return new(Body_0x1206)
	},
	Msg_0x1207: func() Entity {
		return new(Body_0x1207)
	},

	Msg_0x8000: func() Entity {
		return new(Body_0x8000)
	},
	Msg_0x8001: func() Entity {
		return new(Body_0x8001)
	},
	Msg_0x8300: func() Entity {
		return new(Body_0x8300)
	},
	Msg_0x8101: func() Entity {
		return new(Body_0x8101)
	},
	Msg_0x8002: func() Entity {
		return new(Body_0x8002)
	},
	Msg_0x8003: func() Entity {
		return new(Body_0x8003)
	},
	Msg_0x8004: func() Entity {
		return new(Body_0x8004)
	},
	Msg_0x8005: func() Entity {
		return new(Body_0x8005)
	},
	Msg_0x8301: func() Entity {
		return new(Body_0x8301)
	},
	Msg_0x8021: func() Entity {
		return new(Body_0x8021)
	},
	Msg_0x8201: func() Entity {
		return new(Body_0x8201)
	},
	Msg_0x8202: func() Entity {
		return new(Body_0x8202)
	},
	Msg_0x8203: func() Entity {
		return new(Body_0x8203)
	},
	Msg_0x8204: func() Entity {
		return new(Body_0x8204)
	},
	Msg_0x8205: func() Entity {
		return new(Body_0x8205)
	},
	Msg_0x8206: func() Entity {
		return new(Body_0x8206)
	},
	Msg_0x8207: func() Entity {
		return new(Body_0x8207)
	},
}
