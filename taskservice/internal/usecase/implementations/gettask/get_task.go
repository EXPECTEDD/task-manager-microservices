package getuc

import (
	"context"
	"errors"
	"log/slog"
	"taskservice/internal/repository/storage"
	geterr "taskservice/internal/usecase/error/gettask"
	getmodel "taskservice/internal/usecase/models/gettask"
)

type GetTaskUC struct {
	log  *slog.Logger
	stor storage.StorageRepo
}

func NewGetTaskUC(log *slog.Logger, stor storage.StorageRepo) *GetTaskUC {
	return &GetTaskUC{
		log:  log,
		stor: stor,
	}
}

func (g *GetTaskUC) Execute(ctx context.Context, in *getmodel.GetTaskInput) (*getmodel.GetTaskOutput, error) {
	const op = "getuc.Execute"

	log := g.log.With(slog.String("op", op))

	log.Info("starting get task")

	task, err := g.stor.Get(ctx, in.TaskId, in.ProjectId)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Info("task not found")
			return getmodel.NewGetTaskOutput(nil), geterr.ErrTaskNotFound
		}
		log.Warn("cannot get task", slog.String("error", err.Error()))
		return getmodel.NewGetTaskOutput(nil), err
	}

	return getmodel.NewGetTaskOutput(task), nil
}
