package jt808

// 请求周期定位
type T808_0x0110 struct {
}

func (entity *T808_0x0110) MsgID() MsgID {
    return MsgT808_0x0002
}

func (entity *T808_0x0110) Encode() ([]byte, error) {
    return nil, nil
}

func (entity *T808_0x0110) Decode(data []byte) (int, error) {
    return 0, nil
}
