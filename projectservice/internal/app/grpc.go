package app

import (
	"log/slog"
	"projectservice/internal/config"
	grpcserver "projectservice/internal/transport/grpc"
	grpchandler "projectservice/internal/transport/grpc/handler"

	"google.golang.org/grpc"
)

func mustLoadGRPCServer(log *slog.Logger, cfg *config.Config, handl *grpchandler.GRPCHandler) *grpcserver.GRPCServer {
	serv := &grpc.Server{}

	return grpcserver.NewGRPCServer(log, handl, serv, cfg.GrpcConf.Host, cfg.GrpcConf.Port)
}
