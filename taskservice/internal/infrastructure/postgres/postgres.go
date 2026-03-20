package postgres

import (
	"context"
	"database/sql"
	taskdomain "taskservice/internal/domain/task"
	posmapper "taskservice/internal/infrastructure/postgres/mapper"
	"taskservice/internal/repository/storage"
	"time"
)

var (
	invalidId uint32 = 0
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) Save(ctx context.Context, td *taskdomain.TaskDomain) (uint32, error) {
	model := posmapper.TaskDomainToModel(td)

	row := p.db.QueryRowContext(ctx, QuerieCreate, model.ProjectId, model.Description, model.Deadline)

	var id uint32

	err := row.Scan(&id)
	if err != nil {
		return invalidId, err
	}

	return id, nil
}

func (p *Postgres) ChangeDescription(ctx context.Context, taskId uint32, projectId uint32, newDescription string) error {
	res, err := p.db.ExecContext(ctx, QuerieUpdateDescription, newDescription, taskId, projectId)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return storage.ErrTaskNotFound
	}

	return nil
}

func (p *Postgres) ChangeDeadline(ctx context.Context, taskId uint32, projectId uint32, newDeadline time.Time) error {
	res, err := p.db.ExecContext(ctx, QuerieUpdateDeadline, newDeadline, taskId, projectId)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return storage.ErrTaskNotFound
	}

	return nil
}
