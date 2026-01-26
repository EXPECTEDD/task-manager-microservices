package registration

import (
	"context"
	"errors"
	"log/slog"
	userdomain "userservice/internal/domain/user"
	"userservice/internal/repository/hasher"
	storagerepo "userservice/internal/repository/storage"
	regerr "userservice/internal/usecase/errors/registration"
	regmodel "userservice/internal/usecase/models/registration"
)

var (
	invalidId uint32 = 0
)

type RegUserUC struct {
	log *slog.Logger

	storage    storagerepo.StorageRepo
	passHasher hasher.PasswordHasher
}

func NewRegUserUC(log *slog.Logger, storage storagerepo.StorageRepo, passHasher hasher.PasswordHasher) *RegUserUC {
	return &RegUserUC{
		log:        log,
		storage:    storage,
		passHasher: passHasher,
	}
}

func (r *RegUserUC) Execute(ctx context.Context, in *regmodel.RegInput) (*regmodel.RegOutput, error) {
	const op = "registration.Execute"
	log := r.log.With(slog.String("op", op), slog.String("email", in.Email))

	log.Info("user registration started")

	ud, err := r.storage.FindByEmail(ctx, in.Email)
	if err != nil && !errors.Is(err, storagerepo.ErrNoRows) {
		log.Warn("registration stopped", slog.String("error", err.Error()))
		return regmodel.NewRegOutput(false), err
	}
	if ud != nil {
		log.Info("registration stopped, user already exists")
		return regmodel.NewRegOutput(false), regerr.ErrUserAlreadyExists
	}

	hashPass, err := r.passHasher.Hash([]byte(in.Password))
	if err != nil {
		log.Warn("registration stopped, impossible to hash", slog.String("error", err.Error()))
		return regmodel.NewRegOutput(false), err
	}

	ud = userdomain.NewUserDomain(
		invalidId,
		in.FirstName,
		in.MiddleName,
		in.LastName,
		string(hashPass),
		in.Email,
	)

	_, err = r.storage.Save(ctx, ud)
	if err != nil {
		log.Warn("registration stopped, failed to save user", slog.String("error", err.Error()))
		return regmodel.NewRegOutput(false), err
	}

	log.Info("user successfully registered")

	return regmodel.NewRegOutput(true), nil
}
