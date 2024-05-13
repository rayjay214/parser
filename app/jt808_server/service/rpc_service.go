package service

import (
	"context"
	"errors"
	"github.com/rayjay214/parser/app/jt808_server/service/proto"
	"github.com/rayjay214/parser/common"
	"github.com/rayjay214/parser/server"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var gJt808Server *server.Server

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

	storage.SetCmdLog(req.Imei, seqNo, req.TimeId)

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
	case "1", "2":
		_, err := session.OpenShortRecord(req.Seconds)
		if err != nil {
			return &resp, err
		}
		storage.SetCmdLog(req.Imei, 10, req.TimeId) //应答中没有seqno, 用type代替
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
	case "1", "2":
		seqNo, err := session.VorRecordSwitch(req.Switch)
		if err != nil {
			return &resp, err
		}
		storage.SetCmdLog(req.Imei, seqNo, req.TimeId)
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
	default:
		resp.Message = "protocol not supported"
		return &resp, errors.New("protocol not supported")
	}

	param := writer.Bytes()
	seqNo, err := session.SetLocMode(param)
	if err != nil {
		return &resp, err
	}
	storage.SetCmdLogMode(req.Imei, seqNo, req.TimeId, req.Mode)

	return &resp, nil
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
	storage.SetCmdLog(req.Imei, seqNo, req.TimeId)

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

func StartRpc(tcpServer *server.Server) {
	gJt808Server = tcpServer
	lis, err := net.Listen("tcp", ":30051")
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
