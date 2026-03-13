package projectrepository

import "context"

type ProjectRepository interface {
	GetOwnerId(ctx context.Context, projectId uint32) (uint32, error)
}
