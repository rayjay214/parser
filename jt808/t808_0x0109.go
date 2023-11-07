package jt808

// 请求同步时间
type T808_0x0109 struct {
}

func (entity *T808_0x0109) MsgID() MsgID {
    return MsgT808_0x0002
}

func (entity *T808_0x0109) Encode() ([]byte, error) {
    return nil, nil
}

func (entity *T808_0x0109) Decode(data []byte) (int, error) {
    return 0, nil
}
