package storage

import (
	"context"
	taskdomain "taskservice/internal/domain/task"
	"time"
)

type StorageRepo interface {
	Save(ctx context.Context, td *taskdomain.TaskDomain) (uint32, error)
	ChangeDescription(ctx context.Context, taskId uint32, newDescription string) error
	ChangeDeadline(ctx context.Context, taskId uint32, newDeadline time.Time) error
}
