package postgres

import (
	"context"
	"database/sql"
	"errors"
	taskdomain "taskservice/internal/domain/task"
	posmapper "taskservice/internal/infrastructure/postgres/mapper"
	posmodels "taskservice/internal/infrastructure/postgres/models"
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

func (p *Postgres) Delete(ctx context.Context, taskId uint32, projectId uint32) error {
	res, err := p.db.ExecContext(ctx, QuerieDelete, taskId, projectId)
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

func (p *Postgres) GetAll(ctx context.Context, projectId uint32) ([]*taskdomain.TaskDomain, error) {
	rows, err := p.db.QueryContext(ctx, QuerieGetAll, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*posmodels.TaskPosModel{}
	for rows.Next() {
		task := &posmodels.TaskPosModel{}

		err := rows.Scan(
			&task.Id,
			&task.ProjectId,
			&task.Description,
			&task.Deadline,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, storage.ErrTasksNotFound
	}

	return posmapper.TaskModelsToDomains(tasks), nil
}

func (p *Postgres) Get(ctx context.Context, taskId uint32, projectId uint32) (*taskdomain.TaskDomain, error) {
	row := p.db.QueryRowContext(ctx, QuerieGet, projectId, taskId)

	task := &posmodels.TaskPosModel{}
	err := row.Scan(
		&task.Id,
		&task.ProjectId,
		&task.Description,
		&task.Deadline,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrTaskNotFound
		}
		return nil, err
	}

	return posmapper.TaskModelToDomain(task), nil
}
