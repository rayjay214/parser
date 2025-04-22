package service

import (
	"context"
	"errors"
	"github.com/rayjay214/parser/app/jt808_server/service/proto"
	"github.com/rayjay214/parser/protocol/common"
	"github.com/rayjay214/parser/server_base/jt808_base"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var gJt808Server *jt808_base.Server

type deviceService struct {
	proto.UnimplementedDeviceServiceServer
}

func (s *deviceService) SendCmd(ctx context.Context, req *proto.SendCmdRequest) (*proto.SendCmdReply, error) {
	var resp proto.SendCmdReply
	resp.Message = "ok"
	session, ok := gJt808Server.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	seqNo, err := session.SendCmd(req.Content)
	if err != nil {
		return &resp, err
	}

	if req.Protocol == "7" {
		storage.SetCmdLogZZE(req.Imei, req.Content, req.TimeId, req.Protocol)
	} else {
		storage.SetCmdLog(req.Imei, seqNo, req.TimeId, req.Protocol)
	}

	return &resp, nil
}

func (s *deviceService) OpenShortRecord(ctx context.Context, req *proto.OpenShortRecordRequest) (*proto.CommonReply, error) {
	var resp proto.CommonReply
	resp.Message = "ok"
	session, ok := gJt808Server.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	switch req.Protocol {
	case "1", "2", "5":
		_, err := session.OpenShortRecord(req.Seconds)
		if err != nil {
			return &resp, err
		}
		storage.SetCmdLog(req.Imei, 10, req.TimeId, req.Protocol) //应答中没有seqno, 用type代替
	default:
		resp.Message = "protocol not supported"
		return &resp, errors.New("protocol not supported")
	}

	return &resp, nil
}

func (s *deviceService) VorRecordSwitch(ctx context.Context, req *proto.VorRecordSwitchRequest) (*proto.CommonReply, error) {
	var resp proto.CommonReply
	resp.Message = "ok"
	session, ok := gJt808Server.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	switch req.Protocol {
	case "1", "2", "5":
		seqNo, err := session.VorRecordSwitch(req.Switch)
		if err != nil {
			return &resp, err
		}
		storage.SetCmdLog(req.Imei, seqNo, req.TimeId, req.Protocol)
	default:
		resp.Message = "protocol not supported"
		return &resp, errors.New("protocol not supported")
	}

	return &resp, nil
}

func (s *deviceService) SetLocMode(ctx context.Context, req *proto.SetLocModeRequest) (*proto.CommonReply, error) {
	var resp proto.CommonReply
	resp.Message = "ok"
	session, ok := gJt808Server.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	writer := common.NewWriter()

	content := ""

	switch req.Protocol {
	case "1", "2":
		switch req.Mode {
		case "1": //性能
			writer.WriteByte(1)
			writer.WriteUint16(0)
		case "2": //正常
			writer.WriteByte(7)
			writer.WriteUint16(120)
		case "3": //点名
			writer.WriteByte(8)
			writer.WriteUint16(0)
		case "4": //省电
			writer.WriteByte(7)
			writer.WriteUint16(300)
		}
	case "3":
		switch req.Mode {
		case "1":
			content = "TIME#1#"
		case "2":
			content = "TIME#10#"
		case "4":
			content = "TIME#30#"
		}
	case "5", "7":
		switch req.Mode {
		case "1":
			content = "MODE,2,30,1,1,1#"
		case "2":
			content = "MODE,2,120,1,1,1,1#"
		case "4":
			content = "MODE,2,1800,1,1,1#"
		}
	default:
		resp.Message = "protocol not supported"
		return &resp, errors.New("protocol not supported")
	}

	if req.Protocol == "3" || req.Protocol == "5" || req.Protocol == "7" {
		seqNo, err := session.SendCmd(content)
		if err != nil {
			return &resp, err
		}
		storage.SetCmdLogMode(req.Imei, seqNo, req.TimeId, req.Mode, req.Protocol)
		return &resp, nil
	} else {
		param := writer.Bytes()
		seqNo, err := session.SetLocMode(param)
		if err != nil {
			return &resp, err
		}
		storage.SetCmdLogMode(req.Imei, seqNo, req.TimeId, req.Mode, req.Protocol)
		return &resp, nil
	}
}

func (s *deviceService) Locate(ctx context.Context, req *proto.LocateRequest) (*proto.CommonReply, error) {
	var resp proto.CommonReply
	resp.Message = "ok"
	session, ok := gJt808Server.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	seqNo, err := session.Locate()
	if err != nil {
		return &resp, err
	}
	storage.SetCmdLog(req.Imei, seqNo, req.TimeId, req.Protocol)

	return &resp, nil
}

func (s *deviceService) SetShakeValue(ctx context.Context, req *proto.SetShakeValueRequest) (*proto.CommonReply, error) {
	var resp proto.CommonReply
	resp.Message = "ok"
	session, ok := gJt808Server.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	seqNo, err := session.SetShakeValue(req.ShakeValue)
	if err != nil {
		return &resp, err
	}
	storage.SetCmdLogShakeValue(req.Imei, seqNo, req.TimeId, req.ShakeValue)

	return &resp, nil
}

func (s *deviceService) HandelDeviceCtrl(ctx context.Context, req *proto.HandelDeviceCtrlRequest) (*proto.CommonReply, error) {
	var resp proto.CommonReply
	resp.Message = "ok"
	session, ok := gJt808Server.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	seqNo, err := session.DeviceCtrl(byte(req.Cmd))
	if err != nil {
		return &resp, err
	}
	storage.SetCmdLog(req.Imei, seqNo, req.TimeId, req.Protocol)

	return &resp, nil
}

func (s *deviceService) HandelRestart(ctx context.Context, req *proto.HandelRestartRequest) (*proto.CommonReply, error) {
	var resp proto.CommonReply
	resp.Message = "ok"
	session, ok := gJt808Server.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	seqNo, err := session.DeviceRestart()
	if err != nil {
		return &resp, err
	}
	storage.SetCmdLog(req.Imei, seqNo, req.TimeId, req.Protocol)

	return &resp, nil
}

func StartRpc(tcpServer *jt808_base.Server) {
	gJt808Server = tcpServer
	lis, err := net.Listen("tcp", storage.Conf.Grpc.Host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterDeviceServiceServer(server, &deviceService{})

	log.Println("gRPC server is running on :40051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
