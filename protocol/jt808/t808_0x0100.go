package jt808

import (
	"github.com/rayjay214/parser/protocol/common"
)

// 终端注册
type T808_0x0100 struct {
	// 省份
	ProvinceID uint16
	// 城市
	CityID uint16
	// 制造商
	ManufactureID string
	// 终端型号
	Model string
	// 终端ID
	TerminalID string
	// 车牌颜色
	PlateColor byte
	// 车辆标识
	LicenseNo string
}

func (entity *T808_0x0100) MsgID() MsgID {
	return MsgT808_0x0100
}

func (entity *T808_0x0100) Encode() ([]byte, error) {
	writer := common.NewWriter()

	// 写入省份ID
	writer.WriteUint16(entity.ProvinceID)

	// 写入城市ID
	writer.WriteUint16(entity.CityID)

	// 写入制造商
	writer.Write([]byte(entity.ManufactureID), 5)

	// 写入终端型号
	writer.Write([]byte(entity.Model), 20)

	// 写入终端ID
	writer.Write([]byte(entity.TerminalID), 7)

	// 写入车牌颜色
	writer.WriteByte(entity.PlateColor)

	// 写入车辆标识
	if len(entity.LicenseNo) > 0 {
		if err := writer.WriteString(entity.LicenseNo); err != nil {
			return nil, err
		}
	}
	return writer.Bytes(), nil
}

//兼容2011
func (entity *T808_0x0100) Decode(data []byte) (int, error) {
	/*
	   if len(data) < 37 {
	       return 0, errors.ErrInvalidBody
	   }
	*/
	reader := common.NewReader(data)

	// 读取省份ID
	var err error
	entity.ProvinceID, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取城市ID
	entity.CityID, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取制造商
	manufacturer, err := reader.Read(5)
	if err != nil {
		return 0, err
	}
	entity.ManufactureID = common.BytesToString(manufacturer[:])

	if len(data) < 37 {
		// 读取终端型号
		model, err := reader.Read(8)
		if err != nil {
			return 0, err
		}
		entity.Model = common.BytesToString(model[:])

		// 读取终端ID
		terminalID, err := reader.Read(7)
		if err != nil {
			return 0, err
		}
		entity.TerminalID = common.BytesToString(terminalID[:])

		// 读取车牌颜色
		entity.PlateColor, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取车辆标识
		if reader.Len() > 0 {
			licenseNo, err := reader.ReadString()
			if err != nil {
				return 0, err
			}
			entity.LicenseNo = licenseNo
		}
	} else {
		// 读取终端型号
		model, err := reader.Read(20)
		if err != nil {
			return 0, err
		}
		entity.Model = common.BytesToString(model[:])

		// 读取终端ID
		terminalID, err := reader.Read(7)
		if err != nil {
			return 0, err
		}
		entity.TerminalID = common.BytesToString(terminalID[:])

		// 读取车牌颜色
		entity.PlateColor, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取车辆标识
		if reader.Len() > 0 {
			licenseNo, err := reader.ReadString()
			if err != nil {
				return 0, err
			}
			entity.LicenseNo = licenseNo
		}
	}

	return len(data) - reader.Len(), nil
}
