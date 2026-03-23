package deleteuc

import (
	"context"
	"errors"
	"log/slog"
	"taskservice/internal/repository/storage"
	deleteerr "taskservice/internal/usecase/error/deletetask"
	deletemodel "taskservice/internal/usecase/models/deletetask"
)

type DeleteTaskUC struct {
	log  *slog.Logger
	stor storage.StorageRepo
}

func NewDeleteTaskUC(log *slog.Logger, stor storage.StorageRepo) *DeleteTaskUC {
	return &DeleteTaskUC{
		log:  log,
		stor: stor,
	}
}

func (d *DeleteTaskUC) Execute(ctx context.Context, in *deletemodel.DeleteTaskInput) (*deletemodel.DeleteTaskOutput, error) {
	const op = "deletetask.Execute"

	log := d.log.With(slog.String("op", op))

	log.Info("starting deleting task")

	err := d.stor.Delete(ctx, in.TaskId, in.ProjectId)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Info("task not found")
			return deletemodel.NewDeleteTaskOutput(false), deleteerr.ErrTaskNotFound
		}
		log.Warn("cannot delete task", slog.String("error", err.Error()))
		return deletemodel.NewDeleteTaskOutput(false), err
	}

	log.Info("task deletion completed successfully")

	return deletemodel.NewDeleteTaskOutput(true), nil
}
