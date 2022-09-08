package rpcServer

import (
	"fmt"
	rpcLogger "go-demo/grpc/logger"
	"go-demo/grpc/protobuf/app/go_demo"
	"go-demo/grpc/services/app"
	"go-demo/utils/config"
	"go-demo/utils/logger"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func InitGrpc() error {
	switch config.Config.Mode {
	case "release":
	default:
		grpclog.SetLoggerV2(rpcLogger.NewZapLogger(rpcLogger.InitRcpLog()))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Config.RpcPort))
	if err != nil {
		logger.SystemWarnf("failed to listen: %v", err)
		return err
	}
	s := grpc.NewServer(
		grpc.InitialWindowSize(1<<30),
		grpc.InitialConnWindowSize(1<<30),
		grpc.MaxSendMsgSize(4<<30),
		grpc.MaxRecvMsgSize(4<<30),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{PermitWithoutStream: true}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			//MaxConnectionIdle:     0,
			//MaxConnectionAge:      0,
			//MaxConnectionAgeGrace: 0,
			Time:    time.Duration(10) * time.Second,
			Timeout: time.Duration(3) * time.Second,
		}),
	)

	go_demo.RegisterGoDemoActionServer(s, &app.Server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		logger.SystemWarnf("failed to serve: %v", err)
		return err
	}
	return nil
}
