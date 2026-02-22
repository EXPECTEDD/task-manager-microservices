package app

import (
	"log/slog"
	"projectservice/internal/config"
	grpcserver "projectservice/internal/transport/grpc"
	projectservicev1 "projectservice/proto/projectservice"

	"google.golang.org/grpc"
)

func mustLoadGRPCServer(log *slog.Logger, cfg *config.Config, handl projectservicev1.ProjectServiceServer) *grpcserver.GRPCServer {
	serv := &grpc.Server{}

	return grpcserver.NewGRPCServer(log, handl, serv, cfg.GrpcConf.Host, cfg.GrpcConf.Port)
}
