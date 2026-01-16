package registration

import (
	"context"
	"log/slog"
	regmodel "userservice/internal/usecase/models/registration"
)

type RegUC struct {
	log *slog.Logger
}

func NewRegUC(log *slog.Logger) *RegUC {
	return &RegUC{
		log: log,
	}
}

func (r *RegUC) RegUser(ctx context.Context, in *regmodel.RegInput) (*regmodel.RegOutput, error) {
	panic("not implemented")
}
