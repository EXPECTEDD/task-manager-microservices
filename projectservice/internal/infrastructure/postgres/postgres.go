package postgres

import (
	"context"
	"database/sql"
	"errors"
	projectdomain "projectservice/internal/domain/project"
	posmapper "projectservice/internal/infrastructure/postgres/mapper"
	posmodels "projectservice/internal/infrastructure/postgres/models"
	"projectservice/internal/repository/storage"

	"github.com/lib/pq"
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

func (p *Postgres) Save(ctx context.Context, proj *projectdomain.ProjectDomain) (uint32, error) {
	pm := posmapper.DomainToModel(proj)

	row := p.db.QueryRowContext(ctx, QuerieSave, pm.OwnerId, pm.Name)

	var id uint32

	err := row.Scan(
		&id,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return invalidId, storage.ErrAlreadyExists
			}
		}
		return invalidId, err
	}

	return id, nil
}

func (p *Postgres) Delete(ctx context.Context, ownerId uint32, projectId uint32) error {
	res, err := p.db.ExecContext(ctx, QuerieDelete, projectId, ownerId)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (p *Postgres) GetAll(ctx context.Context, ownerId uint32) ([]*projectdomain.ProjectDomain, error) {
	rows, err := p.db.QueryContext(ctx, QuerieGetAll, ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*posmodels.ProjectPosModel
	for rows.Next() {
		project := &posmodels.ProjectPosModel{}

		err := rows.Scan(
			&project.Id,
			&project.OwnerId,
			&project.Name,
			&project.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(projects) == 0 {
		return nil, storage.ErrNotFound
	}

	return posmapper.ModelsToDomain(projects), nil
}

func (p *Postgres) UpdateName(ctx context.Context, ownerId uint32, projectId uint32, newName string) error {
	res, err := p.db.ExecContext(ctx, QuerieUpdateName, newName, ownerId, projectId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return storage.ErrAlreadyExists
			}
		}
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (p *Postgres) GetProject(ctx context.Context, projectId uint32) (*projectdomain.ProjectDomain, error) {
	row := p.db.QueryRowContext(ctx, QuerieGetProject, projectId)

	model := &posmodels.ProjectPosModel{}

	err := row.Scan(
		&model.Id,
		&model.OwnerId,
		&model.Name,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return posmapper.ModelToDomain(model), nil
}
