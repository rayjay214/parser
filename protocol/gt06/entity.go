package gt06

// 消息体
type Entity interface {
	MsgID() MsgID
	Encode() ([]byte, error)
	Decode([]byte) (int, error)
	GetSeqNo() uint16
}
