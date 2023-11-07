package cmpp

// 消息ID枚举
type MsgID uint32

const (
    MSG_CMPP_CONNECT      MsgID = 0x00000001
    MSG_CMPP_CONNECT_RESP MsgID = 0x80000001

    MSG_CMPP_SUBMIT       MsgID = 0x00000004
    MSG_CMPP_SUBMIT_RESP  MsgID = 0x80000004
    MSG_CMPP_DELIVER      MsgID = 0x00000005
    MSG_CMPP_DELIVER_RESP MsgID = 0x80000005
    /*
    	MSG_CMPP_ACTIVE_TEST MsgID = 0x00000008
    	MSG_CMPP_ACTIVE_TEST_RESP MsgID = 0x80000008
    	MSG_CMPP_TERMINATE MsgID = 0x00000002
    	MSG_CMPP_TERMINATE_RESP MsgID = 0x80000002
    */
)

// 消息实体映射
var entityMapper = map[uint32]func() Entity{
    uint32(MSG_CMPP_CONNECT): func() Entity {
        return new(CMPP_CONNECT)
    },
    uint32(MSG_CMPP_CONNECT_RESP): func() Entity {
        return new(CMPP_CONNECT_RESP)
    },
    uint32(MSG_CMPP_SUBMIT): func() Entity {
        return new(CMPP_SUBMIT)
    },
    uint32(MSG_CMPP_SUBMIT_RESP): func() Entity {
        return new(CMPP_SUBMIT_RESP)
    },
    uint32(MSG_CMPP_DELIVER): func() Entity {
        return new(CMPP_DELIVER)
    },
    uint32(MSG_CMPP_DELIVER_RESP): func() Entity {
        return new(CMPP_DELIVER_RESP)
    },
    /*
    	MSG_CMPP_ACTIVE_TEST: func() Entity {
    		return new(CMPP_ACTIVE_TEST)
    	},
    	MSG_CMPP_ACTIVE_TEST_RESP: func() Entity {
    		return new(CMPP_ACTIVE_TEST_RESP)
    	},
    	MSG_CMPP_TERMINATE: func() Entity {
    		return new(CMPP_TERMINATE)
    	},
    	MSG_CMPP_TERMINATE_RESP: func() Entity {
    		return new(CMPP_TERMINATE_RESP)
    	},
    */
}
