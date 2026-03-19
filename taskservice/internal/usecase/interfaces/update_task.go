package interfaces

import (
	"context"
	updatemodel "taskservice/internal/usecase/models/updatetask"
)

type UpdateTaskUsecase interface {
	Execute(ctx context.Context, in *updatemodel.UpdateTaskInput) (*updatemodel.UpdateTaskOutput, error)
}
