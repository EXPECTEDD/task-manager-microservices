package grpcserv

import (
	"fmt"
	"log/slog"
	"net"
	userservicev1 "userservice/proto/userservice"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	log  *slog.Logger
	serv *grpc.Server
	host string
	port uint32
}

func NewGRPCServer(log *slog.Logger, handl userservicev1.UserServiceServer, serv *grpc.Server, host string, port uint32) *GRPCServer {
	userservicev1.RegisterUserServiceServer(serv, handl)
	return &GRPCServer{
		log:  log,
		serv: serv,
		host: host,
		port: port,
	}
}

func (g *GRPCServer) MustStart() {
	const op = "grpcserv.MustStart"
	g.log.Info("starting grpc server", slog.String("op", op), slog.Int("port", int(g.port)))
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", g.host, g.port))
	if err != nil {
		panic("failed listen grpc server: " + err.Error())
	}
	defer l.Close()

	if err := g.serv.Serve(l); err != nil {
		panic("failed serv grpc server: " + err.Error())
	}
}

func (g *GRPCServer) Stop() {
	const op = "grpcserv.Stop"
	g.log.Info("start grpc server shutdown", slog.String("op", op))
	g.serv.GracefulStop()
	g.log.Info("grpc server stopped", slog.String("op", op))
}
