package grpchandler

import (
	"context"
	"log/slog"
	projectservicev1 "projectservice/proto/projectservice"
)

type GRPCHandler struct {
	log *slog.Logger

	projectservicev1.UnimplementedProjectServiceServer
}

func NewGRPCServer(log *slog.Logger) *GRPCHandler {
	return &GRPCHandler{
		log: log,
	}
}

func (h *GRPCHandler) GetOwnerId(ctx context.Context, in *projectservicev1.GetOwnerIdRequest) (*projectservicev1.GetOwnerIdResponse, error) {
	panic("not implemented")
}
