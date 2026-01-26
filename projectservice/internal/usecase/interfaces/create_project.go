package interfaces

import (
	"context"
	createmodel "projectservice/internal/usecase/models/createproject"
)

type CreateProjectUsecase interface {
	Execute(ctx context.Context, in *createmodel.CreateProjectInput) (*createmodel.CreateProjectOutput, error)
}