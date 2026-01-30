package storage

import (
	"context"
	projectdomain "projectservice/internal/domain/project"
)

type StorageRepo interface {
	Save(ctx context.Context, proj *projectdomain.ProjectDomain) error
}
