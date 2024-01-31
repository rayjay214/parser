package service

import (
    "context"
    "errors"
    "github.com/rayjay214/parser/server"
    "github.com/rayjay214/parser/service/proto"
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

    session.SendCmd(req.Content)

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

    session.OpenShortRecord(req.Seconds)

    return &resp, nil
}

func StartRpc(tcpServer *server.Server) {
    gJt808Server = tcpServer
    lis, err := net.Listen("tcp", ":40051")
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
