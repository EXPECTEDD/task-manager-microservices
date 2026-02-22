package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	projectservicev1 "projectservice/proto/projectservice"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	log  *slog.Logger
	serv *grpc.Server
	host string
	port uint32
}

func NewGRPCServer(log *slog.Logger, handl projectservicev1.ProjectServiceServer, serv *grpc.Server, host string, port uint32) *GRPCServer {
	projectservicev1.RegisterProjectServiceServer(serv, handl)

	return &GRPCServer{
		log:  log,
		serv: serv,
		host: host,
		port: port,
	}
}

func (g *GRPCServer) MustStart() {
	const op = "grpcserver.MustStart"

	g.log.Info("starting grpc server", slog.String("op", op), slog.String("host", g.host), slog.Int("port", int(g.port)))
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", g.host, g.port))
	if err != nil {
		panic("failed listen grpc server: " + err.Error())
	}
	defer l.Close()

	if err := g.serv.Serve(l); err != nil {
		panic("failed serv grpc server: " + err.Error())
	}
}

func (g *GRPCServer) Stop(ctx context.Context) {
	const op = "grpcserver"
	g.log.Info("start grpc server shutdown", slog.String("op", op))

	done := make(chan struct{})

	go func() {
		g.serv.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		g.log.Warn("grpc shutdown timeout, forcing stop")
		g.serv.Stop()
	case <-done:
		g.log.Info("grpc server stopped", slog.String("op", op))
	}
}
