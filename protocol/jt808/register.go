package jt808

// 消息ID枚举
type MsgID uint16

const (
	// 终端应答
	MsgT808_0x0001 MsgID = 0x0001
	// 终端心跳
	MsgT808_0x0002 MsgID = 0x0002
	// 假心跳
	MsgT808_0x0f02 MsgID = 0x0f02
	// 终端注销
	MsgT808_0x0003 MsgID = 0x0003
	// 终端注册
	MsgT808_0x0100 MsgID = 0x0100
	// 终端鉴权
	MsgT808_0x0102 MsgID = 0x0102
	// 查询终端参数应答
	MsgT808_0x0104 MsgID = 0x0104
	// 查询终端属性应答
	MsgT808_0x0107 MsgID = 0x0107
	// 终端升级结果通知
	MsgT808_0x0108 MsgID = 0x0108
	// 汇报位置
	MsgT808_0x0200 MsgID = 0x0200
	// 位置信息查询应答
	MsgT808_0x0201 MsgID = 0x0201
	// 事件报告
	MsgT808_0x0301 MsgID = 0x0301
	// 提问答案
	MsgT808_0x0302 MsgID = 0x0302
	// 信息点播/取消
	MsgT808_0x0303 MsgID = 0x0303
	// 车辆控制
	MsgT808_0x0500 MsgID = 0x0500
	// 行驶记录数据上传
	MsgT808_0x0700 MsgID = 0x0700
	// 电子运单上报
	MsgT808_0x0701 MsgID = 0x0701
	// 驾驶员身份信息采集上报
	MsgT808_0x0702 MsgID = 0x0702
	// 定位数据批量上传
	MsgT808_0x0704 MsgID = 0x0704
	// CAN总线数据上传
	MsgT808_0x0705 MsgID = 0x0705
	// 多媒体事件信息上传
	MsgT808_0x0800 MsgID = 0x0800
	// 多媒体数据上传
	MsgT808_0x0801 MsgID = 0x0801
	// 存储多媒体数据检索应答
	MsgT808_0x0802 MsgID = 0x0802
	// 摄像头立即拍摄命令应答
	MsgT808_0x0805 MsgID = 0x0805
	// 平台通用应答
	MsgT808_0x8001 MsgID = 0x8001
	// 补传分包请求
	MsgT808_0x8003 MsgID = 0x8003
	// 终端注册应答
	MsgT808_0x8100 MsgID = 0x8100
	// 设置终端参数
	MsgT808_0x8103 MsgID = 0x8103
	// 查询终端参数
	MsgT808_0x8104 MsgID = 0x8104
	// 终端控制
	MsgT808_0x8105 MsgID = 0x8105
	// 查询指定参数
	MsgT808_0x8106 MsgID = 0x8106
	// 查询终端属性
	MsgT808_0x8107 MsgID = 0x8107
	// 下发终端升级包
	MsgT808_0x8108 MsgID = 0x8108
	// 查询车辆位置
	MsgT808_0x8201 MsgID = 0x8201
	// 临时位置跟踪控制
	MsgT808_0x8202 MsgID = 0x8202
	// 人工确认报警消息
	MsgT808_0x8203 MsgID = 0x8203
	// 文本信息下发
	MsgT808_0x8300 MsgID = 0x8300
	// 事件设置
	MsgT808_0x8301 MsgID = 0x8301
	// 提问下发
	MsgT808_0x8302 MsgID = 0x8302
	// 位置点播菜单设置
	MsgT808_0x8303 MsgID = 0x8303
	// 信息服务
	MsgT808_0x8304 MsgID = 0x8304
	// 电话回拨
	MsgT808_0x8400 MsgID = 0x8400
	// 设置电话本
	MsgT808_0x8401 MsgID = 0x8401
	// 车门控制
	MsgT808_0x8500 MsgID = 0x8500
	// 设置圆形区域
	MsgT808_0x8600 MsgID = 0x8600
	// 删除圆形区域
	MsgT808_0x8601 MsgID = 0x8601
	// 设置矩形区域
	MsgT808_0x8602 MsgID = 0x8602
	// 删除矩形区域
	MsgT808_0x8603 MsgID = 0x8603
	// 设置多边形区域
	MsgT808_0x8604 MsgID = 0x8604
	// 删除多边形区域
	MsgT808_0x8605 MsgID = 0x8605
	// 设置路线
	MsgT808_0x8606 MsgID = 0x8606
	// 删除路线
	MsgT808_0x8607 MsgID = 0x8607
	// 行驶记录数据采集命令
	MsgT808_0x8700 MsgID = 0x8700
	// 行驶记录参数下传命令
	MsgT808_0x8701 MsgID = 0x8701
	// 上报驾驶员身份信息请求
	MsgT808_0x8702 MsgID = 0x8702
	// 多媒体数据上传应答
	MsgT808_0x8800 MsgID = 0x8800
	// 摄像头立即拍摄命令
	MsgT808_0x8801 MsgID = 0x8801
	// 存储多媒体数据检索
	MsgT808_0x8802 MsgID = 0x8802
	// 存储多媒体数据上传命令
	MsgT808_0x8803 MsgID = 0x8803
	// 录音开始命令
	MsgT808_0x8804 MsgID = 0x8804
	// 单条存储多媒体数据检索上传命令
	MsgT808_0x8805 MsgID = 0x8805
	// 数据下行透传
	MsgT808_0x8900 MsgID = 0x8900
	// 数据上行透传
	MsgT808_0x0900 MsgID = 0x0900
	// 数据压缩上报
	MsgT808_0x0901 MsgID = 0x0901
	// 终端 RSA公钥
	MsgT808_0x0A00 MsgID = 0x0a00
	// 平台 RSA公钥
	MsgT808_0x8A00 MsgID = 0x8a00
	// 请求同步时间
	MsgT808_0x0109 MsgID = 0x0109
	// 请求同步时间应答
	MsgT808_0x8109 MsgID = 0x8109
	// 请求周期定位
	MsgT808_0x0110 MsgID = 0x0110
	// 声控开始
	MsgT808_0x0120 MsgID = 0x0120
	// 请求周期定位应答
	MsgT808_0x8110 MsgID = 0x8110
	// 上传模式到服务器
	MsgT808_0x0112 MsgID = 0x0112
	// 请求服务器下发的位置信息
	MsgT808_0x8111 MsgID = 0x8111
	// 服务器下发录音取消
	MsgT808_0x8115 MsgID = 0x8115
	// 下发录音
	MsgT808_0x8116 MsgID = 0x8116
	// 服务器下发录音数据应答
	MsgT808_0x8117 MsgID = 0x8117
	// 服务器下发录音数据应答
	MsgT808_0x8118 MsgID = 0x8118
	// 终端回复录音取消
	MsgT808_0x0115 MsgID = 0x0115
	// 设备应答录音
	MsgT808_0x0116 MsgID = 0x0116
	// 设备上报短录音数据
	MsgT808_0x0117 MsgID = 0x0117
	// 设备上报实时（声控）录音数据
	MsgT808_0x0118 MsgID = 0x0118
	// 设备上报录音数据上传完成通知
	MsgT808_0x0119 MsgID = 0x0119
	// 一键休眠
	MsgT808_0x8135 MsgID = 0x8135
	// 一键唤醒
	MsgT808_0x8145 MsgID = 0x8145
	// 一键重启
	MsgT808_0x8155 MsgID = 0x8155
	// 短信中心号码
	MsgT808_0x0818 MsgID = 0x0818
	// 休眠电池电量更新
	MsgT808_0x0210 MsgID = 0x0210
	// 标准808协议通知
	MsgT808_0x0808 MsgID = 0x0808
	// 请求分发服务器
	MsgT808_0x1004 MsgID = 0x1004
	// 上传ICCID
	MsgT808_0x1107 MsgID = 0x1107
	// 外电电压、电流上报
	MsgT808_0x1005 MsgID = 0x1005
	// 休眠断网开关状态同步
	MsgT808_0x1006 MsgID = 0x1006
	// 时区同步
	MsgT808_0x1007 MsgID = 0x1007
	// 短信应答透传协议
	MsgT808_0x1300 MsgID = 0x1300
	// 恢复出厂设置通知
	MsgT808_0x0106 MsgID = 0x0106
	// 服务器下发短信透传
	MsgT808_0x8113 MsgID = 0x8113
	// 休眠通知
	MsgT808_0x0105 MsgID = 0x0105
	// 休眠通知回复
	MsgT808_0x8125 MsgID = 0x8125
	// 定时开关机
	MsgT808_0x8131 MsgID = 0x8131
	// 上报文本信息
	MsgT808_0x6006 MsgID = 0x6006
	// 上报蓝牙相关模式
	MsgT808_0x0113 MsgID = 0x0113
)

// 消息实体映射
var entityMapper = map[uint16]func() Entity{
	uint16(MsgT808_0x0001): func() Entity {
		return new(T808_0x0001)
	},
	uint16(MsgT808_0x0002): func() Entity {
		return new(T808_0x0002)
	},
	uint16(MsgT808_0x0f02): func() Entity {
		return new(T808_0x0f02)
	},
	uint16(MsgT808_0x0003): func() Entity {
		return new(T808_0x0003)
	},
	uint16(MsgT808_0x0100): func() Entity {
		return new(T808_0x0100)
	},
	uint16(MsgT808_0x0102): func() Entity {
		return new(T808_0x0102)
	},
	uint16(MsgT808_0x0104): func() Entity {
		return new(T808_0x0104)
	},
	uint16(MsgT808_0x0107): func() Entity {
		return new(T808_0x0107)
	},
	uint16(MsgT808_0x0108): func() Entity {
		return new(T808_0x0108)
	},
	uint16(MsgT808_0x0200): func() Entity {
		return new(T808_0x0200)
	},
	uint16(MsgT808_0x0201): func() Entity {
		return new(T808_0x0201)
	},
	uint16(MsgT808_0x0301): func() Entity {
		return new(T808_0x0301)
	},
	uint16(MsgT808_0x0302): func() Entity {
		return new(T808_0x0302)
	},
	uint16(MsgT808_0x0303): func() Entity {
		return new(T808_0x0303)
	},
	uint16(MsgT808_0x0500): func() Entity {
		return new(T808_0x0500)
	},
	uint16(MsgT808_0x0700): func() Entity {
		return new(T808_0x0700)
	},
	uint16(MsgT808_0x0701): func() Entity {
		return new(T808_0x0701)
	},
	uint16(MsgT808_0x0702): func() Entity {
		return new(T808_0x0702)
	},
	uint16(MsgT808_0x0704): func() Entity {
		return new(T808_0x0704)
	},
	uint16(MsgT808_0x0705): func() Entity {
		return new(T808_0x0705)
	},
	uint16(MsgT808_0x0800): func() Entity {
		return new(T808_0x0800)
	},
	uint16(MsgT808_0x0801): func() Entity {
		return new(T808_0x0801)
	},
	uint16(MsgT808_0x0802): func() Entity {
		return new(T808_0x0802)
	},
	uint16(MsgT808_0x0805): func() Entity {
		return new(T808_0x0805)
	},
	uint16(MsgT808_0x8001): func() Entity {
		return new(T808_0x8001)
	},
	uint16(MsgT808_0x8003): func() Entity {
		return new(T808_0x8003)
	},
	uint16(MsgT808_0x8100): func() Entity {
		return new(T808_0x8100)
	},
	uint16(MsgT808_0x8103): func() Entity {
		return new(T808_0x8103)
	},
	uint16(MsgT808_0x8104): func() Entity {
		return new(T808_0x8104)
	},
	uint16(MsgT808_0x8105): func() Entity {
		return new(T808_0x8105)
	},
	uint16(MsgT808_0x8106): func() Entity {
		return new(T808_0x8106)
	},
	uint16(MsgT808_0x8107): func() Entity {
		return new(T808_0x8107)
	},
	uint16(MsgT808_0x8108): func() Entity {
		return new(T808_0x8108)
	},
	uint16(MsgT808_0x8201): func() Entity {
		return new(T808_0x8201)
	},
	uint16(MsgT808_0x8202): func() Entity {
		return new(T808_0x8202)
	},
	uint16(MsgT808_0x8203): func() Entity {
		return new(T808_0x8203)
	},
	uint16(MsgT808_0x8300): func() Entity {
		return new(T808_0x8300)
	},
	uint16(MsgT808_0x8301): func() Entity {
		return new(T808_0x8301)
	},
	uint16(MsgT808_0x8302): func() Entity {
		return new(T808_0x8302)
	},
	uint16(MsgT808_0x8303): func() Entity {
		return new(T808_0x8303)
	},
	uint16(MsgT808_0x8304): func() Entity {
		return new(T808_0x8304)
	},
	uint16(MsgT808_0x8400): func() Entity {
		return new(T808_0x8400)
	},
	uint16(MsgT808_0x8401): func() Entity {
		return new(T808_0x8401)
	},
	uint16(MsgT808_0x8500): func() Entity {
		return new(T808_0x8500)
	},
	uint16(MsgT808_0x8600): func() Entity {
		return new(T808_0x8600)
	},
	uint16(MsgT808_0x8601): func() Entity {
		return new(T808_0x8601)
	},
	uint16(MsgT808_0x8602): func() Entity {
		return new(T808_0x8602)
	},
	uint16(MsgT808_0x8603): func() Entity {
		return new(T808_0x8603)
	},
	uint16(MsgT808_0x8604): func() Entity {
		return new(T808_0x8604)
	},
	uint16(MsgT808_0x8605): func() Entity {
		return new(T808_0x8605)
	},
	uint16(MsgT808_0x8606): func() Entity {
		return new(T808_0x8606)
	},
	uint16(MsgT808_0x8607): func() Entity {
		return new(T808_0x8607)
	},
	uint16(MsgT808_0x8700): func() Entity {
		return new(T808_0x8700)
	},
	uint16(MsgT808_0x8701): func() Entity {
		return new(T808_0x8701)
	},
	uint16(MsgT808_0x8702): func() Entity {
		return new(T808_0x8702)
	},
	uint16(MsgT808_0x8800): func() Entity {
		return new(T808_0x8800)
	},
	uint16(MsgT808_0x8801): func() Entity {
		return new(T808_0x8801)
	},
	uint16(MsgT808_0x8802): func() Entity {
		return new(T808_0x8802)
	},
	uint16(MsgT808_0x8803): func() Entity {
		return new(T808_0x8803)
	},
	uint16(MsgT808_0x8804): func() Entity {
		return new(T808_0x8804)
	},
	uint16(MsgT808_0x8805): func() Entity {
		return new(T808_0x8805)
	},
	uint16(MsgT808_0x8900): func() Entity {
		return new(T808_0x8900)
	},
	uint16(MsgT808_0x0900): func() Entity {
		return new(T808_0x0900)
	},
	uint16(MsgT808_0x0901): func() Entity {
		return new(T808_0x0901)
	},
	uint16(MsgT808_0x0A00): func() Entity {
		return new(T808_0x0A00)
	},
	uint16(MsgT808_0x8A00): func() Entity {
		return new(T808_0x8A00)
	},
	uint16(MsgT808_0x0109): func() Entity {
		return new(T808_0x0109)
	},
	uint16(MsgT808_0x8109): func() Entity {
		return new(T808_0x8109)
	},
	uint16(MsgT808_0x0110): func() Entity {
		return new(T808_0x0110)
	},
	uint16(MsgT808_0x8110): func() Entity {
		return new(T808_0x8110)
	},
	uint16(MsgT808_0x0112): func() Entity {
		return new(T808_0x0112)
	},
	uint16(MsgT808_0x8111): func() Entity {
		return new(T808_0x8111)
	},
	uint16(MsgT808_0x8115): func() Entity {
		return new(T808_0x8115)
	},
	uint16(MsgT808_0x8116): func() Entity {
		return new(T808_0x8116)
	},
	uint16(MsgT808_0x8117): func() Entity {
		return new(T808_0x8117)
	},
	uint16(MsgT808_0x8118): func() Entity {
		return new(T808_0x8118)
	},
	uint16(MsgT808_0x0115): func() Entity {
		return new(T808_0x0115)
	},
	uint16(MsgT808_0x0116): func() Entity {
		return new(T808_0x0116)
	},
	uint16(MsgT808_0x0117): func() Entity {
		return new(T808_0x0117)
	},
	uint16(MsgT808_0x0118): func() Entity {
		return new(T808_0x0118)
	},
	uint16(MsgT808_0x0119): func() Entity {
		return new(T808_0x0119)
	},
	uint16(MsgT808_0x8135): func() Entity {
		return new(T808_0x8135)
	},
	uint16(MsgT808_0x8145): func() Entity {
		return new(T808_0x8145)
	},
	uint16(MsgT808_0x8155): func() Entity {
		return new(T808_0x8155)
	},
	uint16(MsgT808_0x0818): func() Entity {
		return new(T808_0x0818)
	},
	uint16(MsgT808_0x0210): func() Entity {
		return new(T808_0x0210)
	},
	uint16(MsgT808_0x0808): func() Entity {
		return new(T808_0x0808)
	},
	uint16(MsgT808_0x1004): func() Entity {
		return new(T808_0x1004)
	},
	uint16(MsgT808_0x1107): func() Entity {
		return new(T808_0x1107)
	},
	uint16(MsgT808_0x1005): func() Entity {
		return new(T808_0x1005)
	},
	uint16(MsgT808_0x1006): func() Entity {
		return new(T808_0x1006)
	},
	uint16(MsgT808_0x1007): func() Entity {
		return new(T808_0x1007)
	},
	uint16(MsgT808_0x1300): func() Entity {
		return new(T808_0x1300)
	},
	uint16(MsgT808_0x0106): func() Entity {
		return new(T808_0x0106)
	},
	uint16(MsgT808_0x8113): func() Entity {
		return new(T808_0x8113)
	},
	uint16(MsgT808_0x8125): func() Entity {
		return new(T808_0x8125)
	},
	uint16(MsgT808_0x0105): func() Entity {
		return new(T808_0x0105)
	},
	uint16(MsgT808_0x8131): func() Entity {
		return new(T808_0x8131)
	},
	uint16(MsgT808_0x0113): func() Entity {
		return new(T808_0x0113)
	},
	uint16(MsgT808_0x0120): func() Entity {
		return new(T808_0x0120)
	},
	uint16(MsgT808_0x6006): func() Entity {
		return new(T808_0x6006)
	},
}

// 类型注册
func Register(typ uint16, creator func() Entity) {
	entityMapper[typ] = creator
}
