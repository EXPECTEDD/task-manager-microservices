package createproject

import (
	"context"
	"errors"
	"log/slog"
	projectdomain "projectservice/internal/domain/project"
	"projectservice/internal/repository/storage"
	createerr "projectservice/internal/usecase/error/createproject"
	createmodel "projectservice/internal/usecase/models/createproject"
)

type CreateProjectUC struct {
	log *slog.Logger

	stor storage.Storage
}

func NewCreateProjectUC(log *slog.Logger, stor storage.Storage) *CreateProjectUC {
	return &CreateProjectUC{
		log:  log,
		stor: stor,
	}
}

func (c *CreateProjectUC) Execute(ctx context.Context, in *createmodel.CreateProjectInput) (*createmodel.CreateProjectOutput, error) {
	const op = "createproject.Execute"

	log := c.log.With(slog.String("op", op), slog.Int("id", int(in.OwnerId)))

	log.Info("starting create project")

	proj, err := projectdomain.NewProjectDomain(in.OwnerId, in.Name)
	if err != nil {
		log.Info("can't create a project", slog.String("error", err.Error()))
		return &createmodel.CreateProjectOutput{IsCreated: false}, err
	}

	err = c.stor.Save(ctx, proj)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			log.Info("project already exists")
			return &createmodel.CreateProjectOutput{IsCreated: false}, createerr.ErrAlreadyExists
		}
		log.Warn("error save project", slog.String("error", err.Error()))
		return &createmodel.CreateProjectOutput{IsCreated: false}, err
	}

	return &createmodel.CreateProjectOutput{IsCreated: true}, nil
}
