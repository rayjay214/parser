package jt808

import (
	"encoding/json"
	"fmt"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/rayjay214/parser/protocol/jt808/errors"
	"strconv"
)

var MsgIdDesc = map[uint16]string{
	0x0200: "位置信息汇报",
	0x0704: "批量位置信息汇报",
	0x0001: "终端通用应答",
	0x8001: "平台通用应答",
	0x0002: "终端心跳",
	0x0100: "终端注册",
	0x8100: "终端注册应答",
	0x0003: "终端注销",
	0x0102: "终端鉴权",
	0x8103: "设置终端参数",
	0x8104: "查询终端参数",
	0x8106: "查询指定终端参数",
	0x0104: "查询终端参数应答",
	0x8105: "终端控制",
	0x8201: "位置信息查询",
	0x0201: "位置信息查询应答",
	0x8202: "实时追踪",
	0x8203: "人工确认报警信息",
	0x8400: "电话回拨",
	0x8401: "设置电话本",
	0x8107: "查询终端属性",
	0x0107: "查询终端属性应答",
	0x0109: "请求同步时间",
	0x8109: "请求同步时间应答",
	0x8110: "设置周期定位",
	0x0110: "请求周期定位时间",
	0x0112: "同步定位模式",
	0x8111: "发短信请求服务器下发的位置信息",
	0x8116: "下发录音",
	0x0116: "设备应答录音",
	0x0117: "设备上报短录音数据",
	0x8117: "服务器短录音数据应答",
	0x0118: "设备上报实时/声控录音数据",
	0x8118: "服务器实时/声控录音数据应答",
	0x0119: "录音数据上传完成通知",
	0x8115: "服务器下发录音取消",
	0x0115: "终端回复录音取消",
	0x0818: "短信中心号码",
	0x8131: "定时开关机设置",
	0x8135: "一键休眠",
	0x8145: "一键唤醒",
	0x8155: "一键重启",
	0x0210: "休眠电池电量更新",
	0x0808: "标准808协议通知",
	0x1004: "请求分发服务器",
	0x1107: "上传iccid",
	0x1005: "外电电压、电流上报",
	0x1006: "休眠断网开关状态同步",
	0x8300: "文本信息下发",
	0x1300: "短信应答透传协议",
	0x1007: "时区同步指令",
	0x0106: "恢复出厂设置通知",
	0x8113: "服务器下发短信透传",
	0x0108: "休眠唤醒通知",
	0x8108: "平台回复休眠唤醒通知",
	0x0105: "休眠通知",
	0x8125: "平台回复休眠通知",
	0x0120: "声控开始通知",
	0x0113: "蓝牙定位模式上传",
}

// 封包信息
type Packet struct {
	Sum uint16
	Seq uint16
}

// 消息头
type Header struct {
	MsgID       MsgID
	Property    Property
	Version     byte //只有2019有此项
	Imei        uint64
	MsgSerialNo uint16
	Packet      *Packet
	Is2019      bool
}

// 协议编码
func (header *Header) Encode() ([]byte, error) {
	writer := common.NewWriter()

	// 写入消息ID
	writer.WriteUint16(uint16(header.MsgID))

	// 写入消息体属性
	if header.Packet != nil {
		header.Property.enablePacket()
	}
	writer.WriteUint16(uint16(header.Property))

	// 写入终端号码
	writer.Write(common.StringToBCD(strconv.FormatUint(header.Imei, 10), 6))

	// 写入消息流水号
	writer.WriteUint16(header.MsgSerialNo)

	// 写入分包信息
	if header.Property.IsEnablePacket() {
		writer.WriteUint16(header.Packet.Sum)
		writer.WriteUint16(header.Packet.Seq)
	}
	return writer.Bytes(), nil
}

// 协议解码
func (header *Header) Decode(data []byte) error {
	if len(data) < MessageHeaderSize {
		return errors.ErrInvalidHeader
	}
	reader := common.NewReader(data)

	// 读取消息ID
	msgID, err := reader.ReadUint16()
	if err != nil {
		return errors.ErrInvalidHeader
	}

	// 读取消息体属性
	property, err := reader.ReadUint16()
	if err != nil {
		return errors.ErrInvalidHeader
	}

	header.MsgID = MsgID(msgID)
	header.Property = Property(property)

	bodylen := header.Property.GetBodySize()
	if (len(data) - int(bodylen)) == MessageHeaderSize2019 {
		// 读取2019协议版本号
		version, err := reader.ReadByte()
		if err != nil {
			return err
		}
		// 读取终端号码
		temp, err := reader.Read(10)
		if err != nil {
			return errors.ErrInvalidHeader
		}
		imei, err := strconv.ParseUint(common.BcdToString(temp), 10, 64)
		if err != nil {
			return err
		}
		header.Imei = imei
		header.Version = version
		header.Is2019 = true
	} else {
		// 读取终端号码
		temp, err := reader.Read(6)
		if err != nil {
			return errors.ErrInvalidHeader
		}
		imei, err := strconv.ParseUint(common.BcdToString(temp), 10, 64)
		if err != nil {
			return err
		}
		header.Imei = imei
		header.Is2019 = false
	}

	// 读取消息流水号
	serialNo, err := reader.ReadUint16()
	if err != nil {
		return errors.ErrInvalidHeader
	}

	// 读取分包信息
	if Property(property).IsEnablePacket() {
		var packet Packet

		// 读取分包总数
		packet.Sum, err = reader.ReadUint16()
		if err != nil {
			return err
		}

		// 读取分包序列号
		packet.Seq, err = reader.ReadUint16()
		if err != nil {
			return err
		}
		header.Packet = &packet
	}

	header.MsgSerialNo = serialNo
	return nil
}

func (header Header) MarshalJSON() ([]byte, error) {
	type Alias Header

	return json.Marshal(struct {
		Alias
		MsgID   string
		MsgDesc string
	}{
		Alias:   Alias(header),
		MsgID:   "0x" + fmt.Sprintf("%04x", header.MsgID),
		MsgDesc: MsgIdDesc[uint16(header.MsgID)],
	})
}
