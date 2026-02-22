package app

import (
	"log/slog"
	"projectservice/internal/config"
	grpcserver "projectservice/internal/transport/grpc"
	"projectservice/internal/transport/grpc/interceptor"
	projectservicev1 "projectservice/proto/projectservice"

	"google.golang.org/grpc"
)

func mustLoadGRPCServer(log *slog.Logger, cfg *config.Config, handl projectservicev1.ProjectServiceServer) *grpcserver.GRPCServer {
	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.RecoverInterceptor(log),
			interceptor.TimeoutInterceptor(log, cfg.GrpcConf.Timeout),
		),
	)

	return grpcserver.NewGRPCServer(log, handl, serv, cfg.GrpcConf.Host, cfg.GrpcConf.Port)
}
