package postgres

import (
	"context"
	"database/sql"
	"errors"
	userdomain "userservice/internal/domain/user"
	posmapper "userservice/internal/infrastructure/postgres/mapper"
	posmodels "userservice/internal/infrastructure/postgres/models"
	storagerepo "userservice/internal/repository/storage"
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

func (p *Postgres) Save(ctx context.Context, ud *userdomain.UserDomain) (uint32, error) {
	um := posmapper.DomainToModel(ud)

	row := p.db.QueryRowContext(
		ctx,
		QuerySaveUser,
		um.FirstName,
		um.MiddleName,
		um.LastName,
		um.HashPassword,
		um.Email,
	)

	var userId uint32
	err := row.Scan(&userId)
	if err != nil {
		return invalidId, err
	}
	return userId, err
}

func (p *Postgres) FindByEmail(ctx context.Context, email string) (*userdomain.UserDomain, error) {
	row := p.db.QueryRowContext(ctx, QueryFindByEmail, email)

	var um posmodels.UserPosModel

	err := row.Scan(
		&um.Id,
		&um.FirstName,
		&um.MiddleName,
		&um.LastName,
		&um.HashPassword,
		&um.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storagerepo.ErrNoRows
		}
		return nil, err
	}

	return posmapper.ModelToDomain(&um), nil
}
