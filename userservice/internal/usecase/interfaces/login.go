package interfaces

import (
	"context"
	logmodel "userservice/internal/usecase/models/login"
)

type LoginUserUsecase interface {
	Execute(ctx context.Context, in *logmodel.LoginInput) (*logmodel.LoginOutput, error)
}
