package interfaces

import (
	"context"
	regmodel "userservice/internal/usecase/models/registration"
)

type RegistrationUsecase interface {
	RegUser(ctx context.Context, in *regmodel.RegInput) *regmodel.RegOutput
}
