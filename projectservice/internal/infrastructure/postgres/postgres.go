package postgres

import (
	"context"
	"database/sql"
	projectdomain "projectservice/internal/domain/project"
	posmapper "projectservice/internal/infrastructure/postgres/mapper"
	"projectservice/internal/repository/storage"

	"github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) Save(ctx context.Context, proj *projectdomain.ProjectDomain) error {
	pm := posmapper.DomainToModel(proj)

	_, err := p.db.ExecContext(ctx, QuerieSave, pm.OwnerId, pm.Name)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return storage.ErrAlreadyExists
			}
		}
		return err
	}

	return nil
}
