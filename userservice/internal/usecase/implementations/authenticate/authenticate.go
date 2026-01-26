package authenticate

import (
	"context"
	"errors"
	"log/slog"
	"userservice/internal/repository/session"
	autherr "userservice/internal/usecase/errors/authenticate"
	authmodel "userservice/internal/usecase/models/authenticate"
)

var (
	invalidId uint32 = 0
)

type GetUserIDBySessionUC struct {
	log *slog.Logger

	sessionRepo session.SessionRepo
}

func NewGetUserIDBySessionUC(log *slog.Logger, sessionRepo session.SessionRepo) *GetUserIDBySessionUC {
	return &GetUserIDBySessionUC{
		log:         log,
		sessionRepo: sessionRepo,
	}
}

func (a *GetUserIDBySessionUC) Execute(ctx context.Context, in *authmodel.AuthInput) (*authmodel.AuthOutput, error) {
	const op = "authenticate.Execute"
	log := a.log.With(slog.String("op", op))

	log.Info("authenticate session starting")

	userId, err := a.sessionRepo.Get(ctx, in.SessionId)
	if err != nil {
		if errors.Is(err, session.ErrKeyNotFound) {
			log.Info("authenticate stopped: session not found")
			return authmodel.NewAuthOutput(invalidId), autherr.ErrSessionNotFound
		}
		log.Warn("authenticate stopped", slog.String("error", err.Error()))
		return authmodel.NewAuthOutput(invalidId), err
	}

	return authmodel.NewAuthOutput(userId), err
}
