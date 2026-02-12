package storage

import (
	"context"
	projectdomain "projectservice/internal/domain/project"
)

type StorageRepo interface {
	Save(ctx context.Context, proj *projectdomain.ProjectDomain) (uint32, error)
	Delete(ctx context.Context, ownerId uint32, projectId uint32) error
	GetAll(ctx context.Context, ownerId uint32) ([]*projectdomain.ProjectDomain, error)
	UpdateName(ctx context.Context, ownerId uint32, projectId uint32, newName string) error
}
