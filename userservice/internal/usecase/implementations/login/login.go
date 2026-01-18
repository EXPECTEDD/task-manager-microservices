package login

import (
	"context"
	"errors"
	"log/slog"
	"userservice/internal/repository/hasher"
	"userservice/internal/repository/idgenerator"
	"userservice/internal/repository/session"
	storagerepo "userservice/internal/repository/storage"
	logerr "userservice/internal/usecase/errors/login"
	logmodel "userservice/internal/usecase/models/login"
)

var (
	invalidSessionId = "invalid"
)

type LoginUC struct {
	log *slog.Logger

	storage     storagerepo.StorageRepo
	passHasher  hasher.PasswordHasher
	sessionRepo session.SessionRepo
	idgen       idgenerator.IDGenerator
}

func NewLoginUC(
	log *slog.Logger,
	storage storagerepo.StorageRepo,
	passHasher hasher.PasswordHasher,
	sessionRepo session.SessionRepo,
	idgen idgenerator.IDGenerator,
) *LoginUC {
	return &LoginUC{
		log:         log,
		storage:     storage,
		passHasher:  passHasher,
		sessionRepo: sessionRepo,
		idgen:       idgen,
	}
}

func (l *LoginUC) Login(ctx context.Context, in *logmodel.LoginInput) (*logmodel.LoginOutput, error) {
	const op = "login.Login"
	log := l.log.With(slog.String("op", op), slog.String("email", in.Email))

	log.Info("user login started")

	ud, err := l.storage.FindByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, storagerepo.ErrNoRows) {
			log.Info("login stopped: user not found")
			return &logmodel.LoginOutput{}, logerr.ErrUserNotFound
		}
		log.Warn("login stopped", slog.String("error", err.Error()))
		return &logmodel.LoginOutput{}, err
	}

	if err := l.passHasher.ComparePassword([]byte(ud.HashPassword), []byte(in.Password)); err != nil {
		if errors.Is(err, hasher.ErrWrongPassword) {
			log.Info("login stopped: wrong password")
			return &logmodel.LoginOutput{}, logerr.ErrInvalidPassword
		}
		log.Warn("login stopped", slog.String("error", err.Error()))
		return &logmodel.LoginOutput{}, err
	}

	sessionId := l.idgen.New()

	if err := l.sessionRepo.Save(ctx, sessionId, ud.Id); err != nil {
		log.Warn("login stopped: cannot save session")
		return &logmodel.LoginOutput{}, err
	}

	log.Info("user successfully login")

	return logmodel.NewLoginOutput(sessionId, ud.FirstName, ud.MiddleName, ud.LastName), nil
}
