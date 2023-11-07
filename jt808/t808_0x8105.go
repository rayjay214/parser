package jt808

import (
    "encoding/json"
    _ "fmt"
    "github.com/rayjay214/parser/common"
    "github.com/rayjay214/parser/jt808/errors"
)

// 终端控制
type T808_0x8105 struct {
    // 命令字
    Cmd byte
    // 命令参数
    Data string
}

// 参数描述
var paramDesc = map[byte]string{
    3:    "终端关机",
    4:    "终端复位(重启)",
    5:    "终端恢复出厂设置",
    6:    "关闭数据通信",
    7:    "关闭所有无线通信",
    8:    "一键休眠",
    9:    "一键唤醒",
    0x64: "断油电（需结合速度或传感器状态判断）",
    0x65: "恢复油电",
    0x66: "强制断油电（立即断油电，不用判断）",
}

func (entity *T808_0x8105) MsgID() MsgID {
    return MsgT808_0x8105
}

func (entity *T808_0x8105) Encode() ([]byte, error) {
    writer := common.NewWriter()
    writer.WriteByte(byte(entity.Cmd))
    if len(entity.Data) > 0 {
        if err := writer.WriteString(entity.Data); err != nil {
            return nil, err
        }
    }
    return writer.Bytes(), nil
}

func (entity *T808_0x8105) Decode(data []byte) (int, error) {
    if len(data) < 1 {
        return 0, errors.ErrInvalidBody
    }

    reader := common.NewReader(data[1:])
    if reader.Len() > 0 {
        data, err := reader.ReadString()
        if err != nil {
            return 0, err
        }
        entity.Data = data
    }

    entity.Cmd = data[0]
    return len(data) - reader.Len() - 1, nil
}

func (entity T808_0x8105) MarshalJSON() ([]byte, error) {
    type Alias T808_0x8105

    type New8105 struct {
        Alias
        Cmd  string
        Data string
    }

    s := New8105{
        Alias: Alias(entity),
    }

    s.Cmd = paramDesc[entity.Cmd]

    return json.Marshal(s)
}
