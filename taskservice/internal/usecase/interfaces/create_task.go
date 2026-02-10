package interfaces

import (
	"context"
	createmodel "taskservice/internal/usecase/models/createtask"
)

type CreateTaskUsecase interface {
	Execute(ctx context.Context, in *createmodel.CreateTaskInput) (*createmodel.CreateTaskOutput, error)
}
