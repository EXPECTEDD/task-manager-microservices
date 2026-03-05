package grpchandler

import (
	"context"
	"errors"
	"log/slog"
	grpchandlmapper "projectservice/internal/transport/grpc/handler/mapper"
	getowneriderr "projectservice/internal/usecase/error/getownerid"
	"projectservice/internal/usecase/interfaces"
	projectservicev1 "projectservice/proto/projectservice"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	log *slog.Logger

	getOwnerIdUC interfaces.GetOwnerIdUsecase

	projectservicev1.UnimplementedProjectServiceServer
}

func NewGRPCServer(log *slog.Logger, getOwnerIdUC interfaces.GetOwnerIdUsecase) *GRPCHandler {
	return &GRPCHandler{
		log:          log,
		getOwnerIdUC: getOwnerIdUC,
	}
}

func (h *GRPCHandler) GetOwnerId(ctx context.Context, in *projectservicev1.GetOwnerIdRequest) (*projectservicev1.GetOwnerIdResponse, error) {
	const op = "grpchandler.GetOwnerId"

	log := h.log.With(slog.String("op", op), slog.Int("projectId", int(in.ProjectId)))

	log.Info("starting get owner id request")

	if in.ProjectId == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid project id")
	}

	input := grpchandlmapper.GetOwnerIdRequestToInput(in)

	out, err := h.getOwnerIdUC.Execute(ctx, input)
	if err != nil {
		if errors.Is(err, getowneriderr.ErrProjectsNotFound) {
			log.Info("project not found")
			return nil, status.Error(codes.NotFound, "project not found")
		} else {
			log.Warn("cannot get owner id", slog.String("error", err.Error()))
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	log.Info("get owner id request completed successfully")

	return &projectservicev1.GetOwnerIdResponse{OwnerId: out.OwnerId}, nil
}
