package service

import (
	"context"
	"github.com/rayjay214/parser/app/gt06_server/service/proto"
	"github.com/rayjay214/parser/server_base/hl3g_base"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var ht3gServer *hl3g_base.Server

type deviceService struct {
	proto.UnimplementedGt06ServiceServer
}

func (s *deviceService) SendCmd(ctx context.Context, req *proto.SendGt06CmdRequest) (*proto.Gt06CommonReply, error) {
	var resp proto.Gt06CommonReply
	resp.Message = "ok"

	return &resp, nil
}

func StartRpc(tcpServer *hl3g_base.Server) {
	ht3gServer = tcpServer
	lis, err := net.Listen("tcp", storage.Conf.Grpc.Host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterGt06ServiceServer(server, &deviceService{})

	log.Println("gRPC server is running on :40052")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
