package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	userdomain "userservice/internal/domain/user"
	storagerepo "userservice/internal/repository/storage"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPostgres_FindByEmail(t *testing.T) {
	tests := []struct {
		testName string
		email    string

		mockRows *sqlmock.Rows
		mockErr  error

		expUser *userdomain.UserDomain
		expErr  error
	}{
		{
			testName: "Success",
			email:    "gmail@gmail.com",

			mockRows: sqlmock.NewRows([]string{"id", "first_name", "middle_name", "last_name", "hash_password", "email"}).
				AddRow(1,
					"Ivan",
					"Ivanovich",
					"Ivanov",
					"somePass",
					"gmail@gmail.com",
				),
			mockErr: nil,

			expUser: userdomain.NewUserDomain(
				1,
				"Ivan",
				"Ivanovich",
				"Ivanov",
				"somePass",
				"gmail@gmail.com",
			),
			expErr: nil,
		}, {
			testName: "User not found",
			email:    "gmail@gmail.com",
			mockErr:  sql.ErrNoRows,
			expUser:  nil,
			expErr:   storagerepo.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			if tt.mockErr != nil {
				mock.ExpectQuery(regexp.QuoteMeta(QueryFindByEmail)).
					WithArgs(tt.email).
					WillReturnError(tt.mockErr)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(QueryFindByEmail)).
					WithArgs(tt.email).
					WillReturnRows(tt.mockRows)
			}

			repo := NewPostgres(db)
			ud, err := repo.FindByEmail(context.Background(), tt.email)

			assert.ErrorIs(t, tt.expErr, err)
			assert.Equal(t, tt.expUser, ud)
		})
	}
}

func TestPostgres_Save(t *testing.T) {
	tests := []struct {
		testName string
		user     userdomain.UserDomain

		expId  uint32
		expErr error
	}{
		{
			testName: "Success",
			user: *userdomain.NewUserDomain(
				1,
				"Ivan",
				"Ivanovich",
				"Ivanov",
				"somePass",
				"gmail@gmail.com",
			),
			expId:  1,
			expErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta(QuerySaveUser)).
				WithArgs(tt.user.FirstName, tt.user.MiddleName, tt.user.LastName, tt.user.HashPassword, tt.user.Email).
				WillReturnRows(mock.NewRows([]string{"id"}).
					AddRow(1))

			repo := NewPostgres(db)
			id, err := repo.Save(context.Background(), &tt.user)
			assert.ErrorIs(t, tt.expErr, err)
			assert.Equal(t, tt.expId, id)
		})
	}
}
