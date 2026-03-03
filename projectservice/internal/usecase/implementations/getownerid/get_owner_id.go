package getownerid

import (
	"context"
	"errors"
	"log/slog"
	"projectservice/internal/repository/storage"
	getowneriderr "projectservice/internal/usecase/error/getownerid"
	getowneridmodel "projectservice/internal/usecase/models/getownerid"
)

var (
	invalidId uint32 = 0
)

type GetOwnerIdUC struct {
	log *slog.Logger

	stor storage.StorageRepo
}

func NewGetOwnerIdUC(log *slog.Logger, stor storage.StorageRepo) *GetOwnerIdUC {
	return &GetOwnerIdUC{
		log:  log,
		stor: stor,
	}
}

func (g *GetOwnerIdUC) Execute(ctx context.Context, in *getowneridmodel.GetOwnerIdInput) (*getowneridmodel.GetOwnerIdOutput, error) {
	const op = "getownerid.Execute"

	log := g.log.With(slog.String("op", op), slog.Int("projectId", int(in.ProjectId)))

	log.Info("starting get owner id")

	proj, err := g.stor.GetProject(ctx, in.ProjectId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Info("project not found")
			return getowneridmodel.NewGetOwnerIdOutput(0), getowneriderr.ErrProjectsNotFound
		}
		log.Warn("cannot get owner id", slog.String("error", err.Error()))
		return getowneridmodel.NewGetOwnerIdOutput(0), err
	}

	log.Info("owner ID received")

	return getowneridmodel.NewGetOwnerIdOutput(proj.OwnerId), nil
}
