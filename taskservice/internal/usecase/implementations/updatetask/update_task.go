package updateuc

import (
	"context"
	"log/slog"
	"taskservice/internal/repository/storage"
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
	panic("not implemeted")
}
