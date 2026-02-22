package interfaces

import (
	"context"
	getowneridmodel "projectservice/internal/usecase/models/getownerid"
)

type GetOwnerIdUsecase interface {
	Execute(ctx context.Context, in *getowneridmodel.GetOwnerIdInput) (*getowneridmodel.GetOwnerIdOutput, error)
}
