package grpchandler

import (
	"context"
	"errors"
	"log/slog"
	"time"
	autherr "userservice/internal/usecase/errors/authenticate"
	"userservice/internal/usecase/interfaces"
	authmodel "userservice/internal/usecase/models/authenticate"
	userservicev1 "userservice/proto/userservice"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	log     *slog.Logger
	timeout time.Duration
	userservicev1.UnimplementedUserServiceServer

	authUC interfaces.AuthenticateUsecase
}

func NewGRPCHandler(log *slog.Logger, authUC interfaces.AuthenticateUsecase) *GRPCHandler {
	return &GRPCHandler{
		log:    log,
		authUC: authUC,
	}
}

func (g *GRPCHandler) GetIdBySession(ctx context.Context, req *userservicev1.GetIdBySessionRequest) (*userservicev1.GetIdBySessionResponse, error) {
	const op = "grpchandler.GetIdBySession"
	log := g.log.With(slog.String("op", op))

	log.Info("start get id by session request")

	in := authmodel.NewAuthInput(req.SessionId)

	out, err := g.authUC.AuthenticateSession(ctx, in)
	if err != nil {
		if errors.Is(err, autherr.ErrSessionNotFound) {
			log.Info("session not found")
			return nil, status.Error(codes.NotFound, "session not found")
		}
		log.Warn("cannotfailed to get user id", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &userservicev1.GetIdBySessionResponse{
		UserId: out.UserId,
	}, nil
}
