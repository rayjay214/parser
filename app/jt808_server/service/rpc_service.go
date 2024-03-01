package service

import (
	"context"
	"errors"
	proto2 "github.com/rayjay214/parser/app/jt808_server/service/proto"
	"github.com/rayjay214/parser/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var gJt808Server *server.Server

type deviceService struct {
	proto2.UnimplementedDeviceServiceServer
}

func (s *deviceService) SendCmd(ctx context.Context, req *proto2.SendCmdRequest) (*proto2.SendCmdReply, error) {
	var resp proto2.SendCmdReply
	resp.Message = "ok"
	session, ok := gJt808Server.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	session.SendCmd(req.Content)

	return &resp, nil
}

func (s *deviceService) OpenShortRecord(ctx context.Context, req *proto2.OpenShortRecordRequest) (*proto2.CommonReply, error) {
	var resp proto2.CommonReply
	resp.Message = "ok"
	session, ok := gJt808Server.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	session.OpenShortRecord(req.Seconds)

	return &resp, nil
}

func (s *deviceService) VorRecordSwitch(ctx context.Context, req *proto2.VorRecordSwitchRequest) (*proto2.CommonReply, error) {
	var resp proto2.CommonReply
	resp.Message = "ok"
	session, ok := gJt808Server.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	session.VorRecordSwitch(req.Switch)

	return &resp, nil
}

func StartRpc(tcpServer *server.Server) {
	gJt808Server = tcpServer
	lis, err := net.Listen("tcp", ":40051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto2.RegisterDeviceServiceServer(server, &deviceService{})

	log.Println("gRPC server is running on :40051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
