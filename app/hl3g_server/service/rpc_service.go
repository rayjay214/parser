package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/rayjay214/parser/app/hl3g_server/service/proto"
	"github.com/rayjay214/parser/server_base/hl3g_base"
	"github.com/rayjay214/parser/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var ht3gServer *hl3g_base.Server

type deviceService struct {
	proto.UnimplementedHl3GServiceServer
}

func (s *deviceService) SendCmd(ctx context.Context, req *proto.SendHl3GCmdRequest) (*proto.Hl3GCommonReply, error) {
	var resp proto.Hl3GCommonReply
	resp.Message = "ok"

	session, ok := ht3gServer.GetSession(req.Imei)

	if !ok {
		log.Errorf("can't find device %v", req.Imei)
		resp.Message = "can't find device"
		return &resp, errors.New("can't find device")
	}

	imeiStr := fmt.Sprintf("%v", req.Imei)

	seqNo := uint16(0)

	switch req.Proto {
	case "RESET":
		session.CommonReply(imeiStr, req.Proto)
	case "CR":
		session.CommonReply(imeiStr, req.Proto)
	case "POWEROFF":
		session.CommonReply(imeiStr, req.Proto)
	case "VERNO":
		session.CommonReply(imeiStr, req.Proto)
	case "TC":
		session.TcReply(imeiStr, req.Proto, "check123456")
	case "FACTORY":
		session.CommonReply(imeiStr, req.Proto)
	case "UPLOAD":
		var mode string
		if req.Content == "120" {
			mode = "2"
		} else if req.Content == "300" {
			mode = "4"
		}
		storage.SetCmdLogMode(req.Imei, seqNo, req.TimeId, mode, req.Protocol)
		session.UploadReply(imeiStr, req.Proto, req.Content)
	case "MODEWORK":
		mode := "1"
		storage.SetCmdLogMode(req.Imei, seqNo, req.TimeId, mode, req.Protocol)
		session.ModeWorkReply(imeiStr, req.Proto, req.Content)
	case "CLOSEMODE":
		storage.SetCmdLogFakeOnline(req.Imei, req.Content)
		session.CloseModeReply(imeiStr, req.Proto, req.Content)
	}

	storage.SetCmdLog(req.Imei, seqNo, req.TimeId, req.Protocol)

	return &resp, nil
}

func StartRpc(tcpServer *hl3g_base.Server) {
	ht3gServer = tcpServer
	lis, err := net.Listen("tcp", storage.Conf.Grpc.Host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterHl3GServiceServer(server, &deviceService{})

	log.Println("gRPC server is running on :40053")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
