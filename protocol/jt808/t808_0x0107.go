package jt808

import (
	_ "encoding/json"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/rayjay214/parser/protocol/jt808/errors"
)

// 查询终端属性应答
type T808_0x0107 struct {
	// 终端类型
	// bit0 ，0：不适用客运车辆， 1：适用客运车辆；
	// bit1 ，0：不适用危险品车辆， 1：适用危险品车辆；
	// bit2 ，0：不适用普通货运车辆， 1：适用普通货运车辆；
	// bit3 ，0：不适用出租车辆， 1：适用出租车辆；
	// bit6 ，0：不支持硬盘录像， 1：支持硬盘录像；
	// bit7 ，0：一体机， 1：分体机。
	TerminalType1 uint8
	TerminalType2 uint8

	// 制造商
	ManufactureID string
	// 终端型号
	Model string `json:"设备型号"`
	// 终端ID
	TerminalID string
	// SIM卡号
	Iccid string
	// 终端硬件版本
	HardwareVersion string
	// 终端固件版本号
	SoftwareVersion string
	// GNSS模块属性
	GNSSProperty byte
	// 通信模块属性
	COMMProperty byte
}

func (entity *T808_0x0107) MsgID() MsgID {
	return MsgT808_0x0107
}

func (entity *T808_0x0107) Encode() ([]byte, error) {
	writer := common.NewWriter()

	// 写入终端类型
	writer.WriteByte(entity.TerminalType1)
	writer.WriteByte(entity.TerminalType2)

	// 写入终端制造商
	writer.Write([]byte(entity.ManufactureID), 5)

	// 写入终端型号
	writer.Write([]byte(entity.Model), 20)

	// 写入终端ID
	writer.Write([]byte(entity.TerminalID), 7)

	// 写入终端Iccid
	writer.Write(common.StringToBCD(entity.Iccid, 10))

	// 写入终端硬件版本
	hardwareVersion := []byte(entity.HardwareVersion)
	writer.WriteByte(byte(len(hardwareVersion)))
	writer.Write(hardwareVersion)

	// 写入终端固件版本号
	softwareVersion := []byte(entity.SoftwareVersion)
	writer.WriteByte(byte(len(softwareVersion)))
	writer.Write(softwareVersion)

	// 写入GNSS模块属性
	writer.WriteByte(entity.GNSSProperty)

	// 写入通信模块属性
	writer.WriteByte(entity.COMMProperty)
	return writer.Bytes(), nil
}

func (entity *T808_0x0107) Decode(data []byte) (int, error) {
	if len(data) < 28 {
		return 0, errors.ErrInvalidBody
	}
	reader := common.NewReader(data)

	// 读取终端类型
	var err error
	entity.TerminalType1, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	if entity.TerminalType1 == 20 { //兼容2011协议
		entity.Iccid, err = reader.ReadString(20)
		if err != nil {
			return 0, err
		}

		typeLen, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		entity.Model, err = reader.ReadString(int(typeLen))
		if err != nil {
			return 0, err
		}

		versionLen, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		entity.SoftwareVersion, err = reader.ReadString(int(versionLen))
		if err != nil {
			return 0, err
		}
	} else {
		entity.TerminalType2, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取终端制造商
		manufacture, err := reader.Read(5)
		if err != nil {
			return 0, err
		}
		entity.ManufactureID = common.BytesToString(manufacture)

		// 读取终端型号
		model, err := reader.Read(20)
		if err != nil {
			return 0, err
		}
		entity.Model = common.BytesToString(model)

		// 读取终端ID
		terminalID, err := reader.Read(7)
		if err != nil {
			return 0, err
		}
		entity.TerminalID = common.BytesToString(terminalID)

		// 读取Iccid
		temp, err := reader.Read(10)
		if err != nil {
			return 0, err
		}
		entity.Iccid = common.BcdToString(temp)

		// 读取终端硬件版本号长度
		size, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取终端硬件版本号
		temp, err = reader.Read(int(size))
		if err != nil {
			return 0, err
		}
		entity.HardwareVersion = common.BytesToString(temp)

		// 读取终端软件版本号长度
		size, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取终端软件版本号
		temp, err = reader.Read(int(size))
		if err != nil {
			return 0, err
		}
		entity.SoftwareVersion = common.BytesToString(temp[:size])

		// 读取GNSS模块属性
		entity.GNSSProperty, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取通信模块属性
		entity.COMMProperty, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}
	}

	return len(data) - reader.Len(), nil
}
