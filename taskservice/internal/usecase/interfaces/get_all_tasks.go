package interfaces

import (
	"context"
	getallmodel "taskservice/internal/usecase/models/getalltasks"
)

type GetAllTasksUsecase interface {
	Execute(ctx context.Context, in *getallmodel.GetAllTasksInput) (*getallmodel.GetALlTasksOutput, error)
}
