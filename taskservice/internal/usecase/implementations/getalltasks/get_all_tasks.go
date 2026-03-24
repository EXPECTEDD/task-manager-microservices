package getalluc

import (
	"context"
	"errors"
	"log/slog"
	"taskservice/internal/repository/storage"
	getallerr "taskservice/internal/usecase/error/getalltasks"
	getallmodel "taskservice/internal/usecase/models/getalltasks"
)

type GetAllTasksUC struct {
	log  *slog.Logger
	stor storage.StorageRepo
}

func NewGetAllTasksUC(log *slog.Logger, stor storage.StorageRepo) *GetAllTasksUC {
	return &GetAllTasksUC{
		log:  log,
		stor: stor,
	}
}

func (g *GetAllTasksUC) Execute(ctx context.Context, in *getallmodel.GetAllTasksInput) (*getallmodel.GetALlTasksOutput, error) {
	const op = "getalluc.Execute"

	log := g.log.With(slog.String("op", op))

	log.Info("starting get all tasks")

	tasks, err := g.stor.GetAll(ctx, in.ProjectId)
	if err != nil {
		if errors.Is(err, storage.ErrTasksNotFound) {
			log.Info("tasks not found")
			return getallmodel.NewGetAllTasksOutput(nil), getallerr.ErrTasksNotFound
		}
		log.Warn("cannot get all tasks", slog.String("error", err.Error()))
		return getallmodel.NewGetAllTasksOutput(nil), err
	}

	log.Info("get all tasks completed successfully")

	return getallmodel.NewGetAllTasksOutput(tasks), nil
}
