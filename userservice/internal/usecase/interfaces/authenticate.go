package interfaces

import (
	"context"
	authmodel "userservice/internal/usecase/models/authenticate"
)

type AuthenticateUsecase interface {
	AuthenticateSession(ctx context.Context, in *authmodel.AuthInput) (*authmodel.AuthOutput, error)
}
