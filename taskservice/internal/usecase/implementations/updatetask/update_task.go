package updateuc

import (
	"context"
	"errors"
	"log/slog"
	"taskservice/internal/repository/storage"
	updatetaskerr "taskservice/internal/usecase/error/updatetask"
	updatemodel "taskservice/internal/usecase/models/updatetask"
)

type UpdateTaskUC struct {
	log  *slog.Logger
	stor storage.StorageRepo
}

func NewUpdateTaskUC(log *slog.Logger, stor storage.StorageRepo) *UpdateTaskUC {
	return &UpdateTaskUC{
		log:  log,
		stor: stor,
	}
}

func (u *UpdateTaskUC) Execute(ctx context.Context, in *updatemodel.UpdateTaskInput) (*updatemodel.UpdateTaskOutput, error) {
	const op = "updateuc.Execute"

	log := u.log.With(slog.String("op", op))

	log.Info("starting updating task")

	if in.NewDescription != nil {
		err := u.stor.ChangeDescription(ctx, in.TaskId, *in.NewDescription)
		if err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				log.Info("task not found")
				return updatemodel.NewUpdateTaskOutput(false), updatetaskerr.ErrTaskNotFound
			}
			log.Warn("cannot update description", slog.String("error", err.Error()))
			return updatemodel.NewUpdateTaskOutput(false), err
		}
	}
	if in.NewDeadline != nil {
		err := u.stor.ChangeDeadline(ctx, in.TaskId, *in.NewDeadline)
		if err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				log.Info("task not found")
				return updatemodel.NewUpdateTaskOutput(false), updatetaskerr.ErrTaskNotFound
			}
			log.Warn("cannot update deadline", slog.String("error", err.Error()))
			return updatemodel.NewUpdateTaskOutput(false), err
		}
	}

	log.Info("update successful")

	return updatemodel.NewUpdateTaskOutput(true), nil
}
