package interfaces

import (
	"context"
	regmodel "userservice/internal/usecase/models/registration"
)

type RegisterUserUsecase interface {
	Execute(ctx context.Context, in *regmodel.RegInput) (*regmodel.RegOutput, error)
}
