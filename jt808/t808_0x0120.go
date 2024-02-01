package jt808

// 声控开始
type T808_0x0120 struct{}

func (entity *T808_0x0120) MsgID() MsgID {
    return MsgT808_0x0120
}

func (entity *T808_0x0120) Encode() ([]byte, error) {
    return nil, nil
}

func (entity *T808_0x0120) Decode(data []byte) (int, error) {
    return 0, nil
}
