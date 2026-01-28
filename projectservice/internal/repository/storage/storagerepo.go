package storage

import "context"

type StorageRepo interface {
	Save(ctx context.Context)
}
