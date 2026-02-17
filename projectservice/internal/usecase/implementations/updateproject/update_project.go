package updateproject

import (
	"context"
	"errors"
	"log/slog"
	"projectservice/internal/repository/storage"
	updateerr "projectservice/internal/usecase/error/updateproject"
	updatemodel "projectservice/internal/usecase/models/updateproject"
)

type UpdateProjectUC struct {
	log *slog.Logger

	stor storage.StorageRepo
}

func NewUpdateProjectUC(log *slog.Logger, stor storage.StorageRepo) *UpdateProjectUC {
	return &UpdateProjectUC{
		log:  log,
		stor: stor,
	}
}

func (u *UpdateProjectUC) Execute(ctx context.Context, in *updatemodel.UpdateProjectInput) (*updatemodel.UpdateProjectOutput, error) {
	const op = "updateproject.Execute"

	log := u.log.With(slog.String("op", op), slog.Int("ownerId", int(in.OwnerId)), slog.Int("projectId", int(in.ProjectId)))

	log.Info("starting update project")

	rNewName := []rune(*in.NewName)
	if len(rNewName) > 255 || len(rNewName) == 0 {
		return updatemodel.NewUpdateProjectOutput(false), updateerr.ErrInvalidName
	}

	if in.NewName != nil {
		err := u.stor.UpdateName(ctx, in.OwnerId, in.ProjectId, *in.NewName)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return updatemodel.NewUpdateProjectOutput(false), updateerr.ErrProjectNotFound
			} else if errors.Is(err, storage.ErrAlreadyExists) {
				return updatemodel.NewUpdateProjectOutput(false), updateerr.ErrProjectNameAlreadyExists
			}
		}
	}

	log.Info("update project completed successfully")

	return updatemodel.NewUpdateProjectOutput(true), nil
}
