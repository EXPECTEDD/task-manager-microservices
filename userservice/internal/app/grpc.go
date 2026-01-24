package app

import (
	"log/slog"
	"userservice/internal/config"
	grpcserv "userservice/internal/transport/grpc"
	"userservice/internal/transport/grpc/interceptor"
	userservicev1 "userservice/proto/userservice"

	"google.golang.org/grpc"
)

func mustLoadGRPCServer(cfg *config.Config, log *slog.Logger, handl userservicev1.UserServiceServer) *grpcserv.GRPCServer {
	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.RecoverInterceptor(log),
		),
	)

	return grpcserv.NewGRPCServer(log, cfg.GrpcConf.Port, handl, serv)
}
