package jt808

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "golang.org/x/text/encoding/simplifiedchinese"
    "golang.org/x/text/transform"
    "io/ioutil"
    "github.com/rayjay214/parser/common"
    "github.com/rayjay214/parser/jt808/errors"
)

var ParamIdDesc = map[uint32]string{
    0x0001: "心跳发送间隔",
    0x0002: "TCP消息应答超时时间",
    0x0003: "TCP消息重传次数",
    0x0010: "主服务器APN",
    0x0011: "主服务器无线通信拨号用户名",
    0x0012: "主服务器无线通信拨号密码",
    0x0013: "主服务器地址，IP或域名",
    0x0018: "服务器TCP端口",
    0x001E: "分发服务器下发的服务器地址，IP或域名",
    0x001F: "分发服务器下发的服务器TCP端口",
    0x0020: "位置汇报策略", //enum
    0x0021: "位置汇报方案", //enum
    0x0027: "休眠时汇报时间间隔",
    0x0028: "紧急报警时汇报时间间隔",
    0x0029: "缺省时间汇报间隔",
    0x002C: "缺省距离汇报间隔",
    0x002D: "驾驶员未登录汇报距离间隔",
    0x002E: "休眠时汇报距离间隔",
    0x002F: "紧急报警时汇报距离间隔",
    0x0030: "拐点补传角度",
    0x0031: "电子围栏半径",
    0x0040: "监控平台电话号码",
    0x0041: "复位电话号码",
    0x0042: "恢复出厂设置电话号码",
    0x0048: "监听电话号码",
    0x0050: "报警屏蔽字",
    0x0055: "最高速度",
    0x0056: "超速持续时间",
    0x0057: "连续驾驶时间门限",
    0x0058: "当天累计驾驶时间门限",
    0x0059: "最小休息时间",
    0x005A: "最长停车时间",
    0x005C: "疲劳驾驶预警差值",
    0x0061: "声控录音开关",
    0x0062: "实时定位上传模式",
    0x0063: "设备当前语言",
    0x0075: "油电状态",
    0x0076: "GPS定位间隔",
    0x0077: "设备灯控制",
    0x0078: "设置震动报警值",
    0x0079: "在速度超过一定值，认为运动上传数据",
    0x0080: "设置强烈震动报警值",
    0x0081: "车辆所在的省域ID",
    0x0082: "车辆所在的市域ID",
    0x0083: "机动车号牌",
    0x0084: "车牌颜色",
    0x0085: "防拆报警开关",
    0x0086: "置中心号码",
    0x0087: "设置A6S_E模式三的静止时的连接间隔",
    0x0091: "低电报警开关",
    0x0092: "远程升级开关",
    0x0095: "定位数据上传",
    0x0096: "LOG上传开关",
    0x0097: "关机报警开关",
    0x0098: "设置终端定位模式",
    0x0099: "监听开关",
    0xF100: "防拆报警开关",
    0xF101: "设置中心号码",
    0xF102: "低电报警开关",
    0xF103: "远程升级开关",
    0xF104: "定位数据上传设置",
    0xF105: "LOG上传开关",
    0xF106: "关机报警",
    0xF107: "定位模式",
    0xF108: "监听开关",
    0xF109: "油电状态",
    0xF110: "GPS定位间隔",
    0xF111: "设备灯控制",
    0xF112: "设置震动报警值",
    0xF113: "在速度超过一定值，认为运动上传数据",
    0xF114: "声控录音开关",
    0xF115: "实时定位上传模式",
    0xF116: "设备当前语言",
    0xF117: "休眠时间间隔",
    0xF118: "设备剩余电池电量",
    0xF119: "在一定时间内达到一定次数(n)判定运动",
    0xF121: "设置终端定位模式",
    0xF122: "声音提示开关",
    0xF130: "从服务器地址",
    0xF131: "从服务器TCP端口",
    0xF140: "设置设备型号",
    0xF141: "设置休眠后动作",
    0xF142: "设置时区",
    0xF145: "MIC开关",
    0xF146: "定时数据延时开关",
    0xF148: "FOTA相关升级包的下载url",
    0xF149: "设置各种报警触发时电话报警次数",
    0xF150: "设置不同报警类型触发时的报警通知方式",
    0xF151: "设防及震动报警参数设置",
    0xF152: "设防模式设置",
    0xF154: "唤醒状态下静止时是否上报定位数据",
    0xF155: "在定到位的状态下，是否带上基站及wifi等附带信息",
    0xF156: "设置当前协议",
    0xF157: "ACC上电报警开关",
    0xF158: "设置LOG上报等级",
    0xF165: "设置蓝牙模式",
}

var ParamIdParser = map[uint32]func(data []byte) interface{}{
    0x0001: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0002: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0003: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0010: func(data []byte) interface{} {
        return string(data)
    },
    0x0011: func(data []byte) interface{} {
        return string(data)
    },
    0x0012: func(data []byte) interface{} {
        return string(data)
    },
    0x0013: func(data []byte) interface{} {
        return string(data)
    },
    0x0018: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x001E: func(data []byte) interface{} {
        return string(data)
    },
    0x001F: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0020: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0021: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0027: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0028: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0029: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x002C: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x002D: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x002E: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x002F: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0030: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0031: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0040: func(data []byte) interface{} {
        return string(data)
    },
    0x0041: func(data []byte) interface{} {
        return string(data)
    },
    0x0042: func(data []byte) interface{} {
        return string(data)
    },
    0x0048: func(data []byte) interface{} {
        return string(data)
    },
    0x0050: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0055: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0056: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0057: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0058: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0059: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x005A: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x005C: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0061: func(data []byte) interface{} {
        return data[0]
    },
    0x0062: func(data []byte) interface{} {
        m := make(map[string]interface{})
        switch data[0] {
        case 0:
            m["实时定位上传模式"] = "定时上传"
        case 1:
            m["实时定位上传模式"] = "定距上传"
        }
        m["时间间隔/距离间隔"] = binary.BigEndian.Uint16(data[1:])
        return m
    },
    0x0063: func(data []byte) interface{} {
        m := make(map[string]interface{})
        switch data[0] {
        case 0:
            m["语言"] = "英文"
        case 1:
            m["语言"] = "中文"
        }
        return m
    },
    0x0075: func(data []byte) interface{} {
        flag := binary.BigEndian.Uint32(data)
        m := make(map[string]interface{})
        switch flag {
        case 0:
            m["油电状态"] = "断油电"
        case 1:
            m["油电状态"] = "恢复油电"
        }
        return m
    },
    0x0076: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0077: func(data []byte) interface{} {
        m := make(map[string]interface{})
        switch data[0] {
        case 0:
            m["设备灯状态"] = "灭"
        case 1:
            m["设备灯状态"] = "亮"
        }
        return m
    },
    0x0078: func(data []byte) interface{} {
        return data[0]
    },
    0x0079: func(data []byte) interface{} {
        return data[0]
    },
    0x0080: func(data []byte) interface{} {
        return data[0]
    },
    0x0081: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0082: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0083: func(data []byte) interface{} {
        return string(data)
    },
    0x0084: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0x0085: func(data []byte) interface{} {
        return data[0]
    },
    0x0086: func(data []byte) interface{} {
        return string(data)
    },
    0x0087: func(data []byte) interface{} {
        return binary.BigEndian.Uint16(data)
    },
    0x0091: func(data []byte) interface{} {
        return data[0]
    },
    0x0092: func(data []byte) interface{} {
        return data[0]
    },
    0x0095: func(data []byte) interface{} {
        return data[0]
    },
    0x0096: func(data []byte) interface{} {
        return data[0]
    },
    0x0097: func(data []byte) interface{} {
        return data[0]
    },
    0x0098: func(data []byte) interface{} {
        m := make(map[string]interface{})
        mode := data[0]
        switch mode {
        case 1:
            m["定位模式"] = "常用模式"
        case 2:
            m["定位模式"] = "周期定位模式"
        case 3:
            m["定位模式"] = "飞行模式"
        case 4:
            m["定位模式"] = "超级省电模式"
        }

        m["短连接周期"] = binary.BigEndian.Uint16(data[1:3])
        m["是否开启飞行模式"] = data[3]
        m["追车模式结束后再维持多久连接"] = data[4]

        return m
    },
    0x0099: func(data []byte) interface{} {
        return data[0]
    },
    0xF100: func(data []byte) interface{} {
        return data[0]
    },
    0xF101: func(data []byte) interface{} {
        return string(data)
    },
    0xF102: func(data []byte) interface{} {
        return data[0]
    },
    0xF103: func(data []byte) interface{} {
        return data[0]
    },
    0xF104: func(data []byte) interface{} {
        return data[0]
    },
    0xF105: func(data []byte) interface{} {
        return data[0]
    },
    0xF106: func(data []byte) interface{} {
        return data[0]
    },
    0xF107: func(data []byte) interface{} {
        m := make(map[string]interface{})
        mode := data[0]
        switch mode {
        case 1:
            m["定位模式"] = "常用模式"
        case 2:
            m["定位模式"] = "周期定位模式"
        case 3:
            m["定位模式"] = "飞行模式"
        case 4:
            m["定位模式"] = "超级省电模式"
        }

        m["短连接周期"] = binary.BigEndian.Uint16(data[1:3])
        m["是否开启飞行模式"] = data[3]
        m["追车模式结束后再维持多久连接"] = data[4]

        return m
    },
    0xF108: func(data []byte) interface{} {
        return data[0]
    },
    0xF109: func(data []byte) interface{} {
        flag := binary.BigEndian.Uint32(data)
        m := make(map[string]interface{})
        switch flag {
        case 0:
            m["油电状态"] = "断开油电"
        case 1:
            m["油电状态"] = "恢复油电"
        }
        return m
    },
    0xF110: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0xF111: func(data []byte) interface{} {
        return data[0]
    },
    0xF112: func(data []byte) interface{} {
        return data[0]
    },
    0xF113: func(data []byte) interface{} {
        return data[0]
    },
    0xF114: func(data []byte) interface{} {
        return data[0]
    },
    0xF115: func(data []byte) interface{} {
        m := make(map[string]interface{})
        switch data[0] {
        case 0:
            m["实时定位上传模式"] = "定时上传"
        case 1:
            m["实时定位上传模式"] = "定距上传"
        }
        m["时间间隔/距离间隔"] = binary.BigEndian.Uint16(data[1:])
        return m
    },
    0xF116: func(data []byte) interface{} {
        m := make(map[string]interface{})
        switch data[0] {
        case 0:
            m["语言"] = "英文"
        case 1:
            m["语言"] = "中文"
        }
        return m
    },
    0xF117: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0xF118: func(data []byte) interface{} {
        return data[0]
    },
    0xF119: func(data []byte) interface{} {
        return data[0]
    },
    0xF121: func(data []byte) interface{} {
        m := make(map[string]interface{})
        mode := data[0]
        switch mode {
        case 1:
            m["模式"] = "常用模式"
        case 2:
            m["模式"] = "周期定位模式"
        case 3:
            m["模式"] = "飞行模式"
        case 4:
            m["模式"] = "超级省电模式"
        case 5:
            m["模式"] = "智能模式"
        case 6:
            m["模式"] = "待机模式"
        case 7:
            m["模式"] = "省电模式"
        case 8:
            m["模式"] = "点名模式"
        }
        if mode == 7 {
            m["间隔"] = fmt.Sprintf("%d秒", binary.BigEndian.Uint16(data[1:3]))
        } else {
            m["时间"] = fmt.Sprintf("%02x:%02x", data[1:2], data[2:3])
        }
        return m
    },
    0xF122: func(data []byte) interface{} {
        return data[0]
    },
    0xF130: func(data []byte) interface{} {
        return string(data)
    },
    0xF131: func(data []byte) interface{} {
        return binary.BigEndian.Uint32(data)
    },
    0xF140: func(data []byte) interface{} {
        return string(data)
    },
    0xF141: func(data []byte) interface{} {
        m := make(map[string]interface{})
        switch data[0] {
        case 0:
            m["休眠后的动作"] = "发送心跳"
        case 1:
            m["休眠后的动作"] = "断网"
        case 2:
            m["休眠后的动作"] = "发送定位"
        }
        return m
    },
    0xF142: func(data []byte) interface{} {
        t := data[0]

        timezone := int(t) / 8
        if common.GetBit(int(t), 0) != 0 {
            timezone = timezone * -1
        } else {
            timezone = timezone
        }
        return timezone
    },
    0xF145: func(data []byte) interface{} {
        return data[0]
    },
    0xF146: func(data []byte) interface{} {
        return data[0]
    },
    0xF150: func(data []byte) interface{} {
        m := make(map[string]interface{})
        switch data[0] {
        case 1:
            m["报警类型"] = "震动报警"
        case 2:
            m["报警类型"] = "断电报警"
        case 3:
            m["报警类型"] = "ACC接通报警"
        }
        switch data[1] {
        case 0:
            m["报警通知方式"] = "不通知"
        case 1:
            m["报警通知方式"] = "短信通知"
        case 2:
            m["报警通知方式"] = "电话通知"
        case 3:
            m["报警通知方式"] = "短信加电话通知"
        }
        return m
    },
    0xF151: func(data []byte) interface{} {
        m := make(map[string]interface{})
        switch data[0] {
        case 0:
            m["开关"] = "关闭"
        case 1:
            m["开关"] = "开启"
        }
        m["自动设防延时时间"] = data[1]
        m["震动检测时间"] = binary.BigEndian.Uint32(data[2:4])
        m["自动设防震动报警延时"] = binary.BigEndian.Uint32(data[4:6])
        m["震动报警间隔"] = binary.BigEndian.Uint32(data[6:])
        return m
    },
    0xF152: func(data []byte) interface{} {
        m := make(map[string]interface{})
        switch data[0] {
        case 0:
            m["设防模式"] = "发送心跳"
        case 1:
            m["设防模式"] = "手动设防"
        }
        return m
    },
    0xF154: func(data []byte) interface{} {
        return data[0]
    },
    0xF155: func(data []byte) interface{} {
        return data[0]
    },
    0xF156: func(data []byte) interface{} {
        m := make(map[string]interface{})
        switch data[0] {
        case 0:
            m["协议"] = "JT808/2011"
        case 1:
            m["协议"] = "JT808/2013"
        case 2:
            m["协议"] = "JT808/2019"
        case 3:
            m["协议"] = "康凯斯"
        }
        return m
    },
    0xF157: func(data []byte) interface{} {
        return data[0]
    },
    0xF158: func(data []byte) interface{} {
        m := make(map[string]interface{})
        switch data[0] {
        case 0:
            m["上报等级"] = "关闭"
        case 1:
            m["上报等级"] = "INFO"
        case 2:
            m["上报等级"] = "WARNING"
        case 3:
            m["上报等级"] = "ERROR"
        case 4:
            m["上报等级"] = "FATAL"
        }
        return m
    },
    0xF165: func(data []byte) interface{} {
        m := make(map[string]interface{})
        mode := data[0]
        switch mode {
        case 15:
            m["定位模式"] = "鹅卵石自由模式"
        case 16:
            m["定位模式"] = "鹅卵石待机模式"
        case 17:
            m["定位模式"] = "鹅卵石定时模式"
        case 18:
            m["定位模式"] = "鹅卵石实时模式"
        }
        binary.BigEndian.Uint16(data[1:3])
        m["是否定时上报"] = data[1]
        m["是否保持公网"] = data[2]
        m["是否等待时间到来才上报"] = data[3]
        m["是否wifi优先"] = data[4]

        m["公网上传间隔时间"] = binary.BigEndian.Uint32(data[5:9])
        m["蓝牙上传间隔时间"] = binary.BigEndian.Uint32(data[9:13])

        return m
    },
}

// 终端参数
type Param struct {
    id         uint32
    serialized []byte
}

// 参数ID
func (param Param) ID() uint32 {
    return param.id
}

// 设为Byte
func (param Param) SetByte(id uint32, b byte) Param {
    param.id = id
    param.serialized = []byte{b}
    return param
}

// 设为Bytes
func (param Param) SetBytes(id uint32, b []byte) Param {
    param.id = id
    buffer := make([]byte, len(b))
    copy(buffer, b)
    param.serialized = buffer
    return param
}

// 设为Uint16
func (param Param) SetUint16(id uint32, n uint16) Param {
    param.id = id
    var buffer [2]byte
    binary.BigEndian.PutUint16(buffer[:], n)
    param.serialized = buffer[:]
    return param
}

// 设为Uint32
func (param Param) SetUint32(id uint32, n uint32) Param {
    param.id = id
    var buffer [4]byte
    binary.BigEndian.PutUint32(buffer[:], n)
    param.serialized = buffer[:]
    return param
}

// 设为字符串
func (param Param) SetString(id uint32, s string) Param {
    if len(s) == 0 {
        return param.SetBytes(id, nil)
    }
    data, _ := ioutil.ReadAll(transform.NewReader(
        bytes.NewReader([]byte(s)), simplifiedchinese.GB18030.NewEncoder()))
    return param.SetBytes(id, data)
}

// 读取Byte
func (param Param) GetByte() (byte, error) {
    if len(param.serialized) < 1 {
        return 0, errors.ErrInvalidBody
    }
    return param.serialized[0], nil
}

// 读取Bytes
func (param Param) GetBytes() ([]byte, error) {
    data := make([]byte, len(param.serialized))
    copy(data, param.serialized)
    return data, nil
}

// 读取Uint16
func (param Param) GetUint16() (uint16, error) {
    if len(param.serialized) < 2 {
        return 0, errors.ErrInvalidBody
    }
    return binary.BigEndian.Uint16(param.serialized[:2]), nil
}

// 读取Uint32
func (param Param) GetUint32() (uint32, error) {
    if len(param.serialized) < 4 {
        return 0, errors.ErrInvalidBody
    }
    return binary.BigEndian.Uint32(param.serialized[:4]), nil
}

// 读取字符串
func (param Param) GetString() (string, error) {
    data, err := ioutil.ReadAll(transform.NewReader(
        bytes.NewReader(param.serialized), simplifiedchinese.GB18030.NewDecoder()))
    if err != nil {
        return "", err
    }
    return string(data), nil
}
