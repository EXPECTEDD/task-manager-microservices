package interfaces

import (
	"context"
	authmodel "userservice/internal/usecase/models/authenticate"
)

type GetUserIDBySessionUsecase interface {
	Execute(ctx context.Context, in *authmodel.AuthInput) (*authmodel.AuthOutput, error)
}
