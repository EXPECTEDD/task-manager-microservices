package interfaces

import (
	"context"
	getmodel "taskservice/internal/usecase/models/gettask"
)

type GetTaskUsecase interface {
	Execute(ctx context.Context, in *getmodel.GetTaskInput) (*getmodel.GetTaskOutput, error)
}
