package service

import (
	"context"
	"errors"
	"github.com/rayjay214/parser/app/gt06_server/service/proto"
	"github.com/rayjay214/parser/server_base/gt06_base"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var gt06Server *gt06_base.Server

type deviceService struct {
	proto.UnimplementedDeviceServiceServer
}

func (s *deviceService) SendCmd(ctx context.Context, req *proto.SendCmdRequest) (*proto.SendCmdReply, error) {
	var resp proto.SendCmdReply
	resp.Message = "ok"
	session, ok := gt06Server.GetSession(req.Imei)

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

func StartRpc(tcpServer *gt06_base.Server) {
	gt06Server = tcpServer
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
