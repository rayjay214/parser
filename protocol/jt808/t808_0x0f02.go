package jt808

// 终端心跳
type T808_0x0f02 struct {
}

func (entity *T808_0x0f02) MsgID() MsgID {
	return MsgT808_0x0f02
}

func (entity *T808_0x0f02) Encode() ([]byte, error) {
	return nil, nil
}

func (entity *T808_0x0f02) Decode(data []byte) (int, error) {
	return 0, nil
}
