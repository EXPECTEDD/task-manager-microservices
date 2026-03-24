package interfaces

import (
	"context"
	deletemodel "taskservice/internal/usecase/models/deletetask"
)

type DeleteTaskUsecase interface {
	Execute(ctx context.Context, in *deletemodel.DeleteTaskInput) (*deletemodel.DeleteTaskOutput, error)
}
