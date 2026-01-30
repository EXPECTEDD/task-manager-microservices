package postgres

import (
	"context"
	projectdomain "projectservice/internal/domain/project"
	"projectservice/internal/repository/storage"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgres_Save(t *testing.T) {
	tests := []struct {
		testName string

		proj *projectdomain.ProjectDomain

		returnLastId       int64
		returnRowsAffected int64
		returnErr          error
		expErr             error
	}{
		{
			testName: "Success",

			proj: &projectdomain.ProjectDomain{OwnerId: 1, Name: "Name"},

			returnLastId:       1,
			returnRowsAffected: 1,
			returnErr:          nil,
			expErr:             nil,
		}, {
			testName: "Already exists",

			proj: &projectdomain.ProjectDomain{OwnerId: 1, Name: "Name"},

			returnLastId:       1,
			returnRowsAffected: 1,
			returnErr:          &pq.Error{Code: "23505"},
			expErr:             storage.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			mock.ExpectExec(regexp.QuoteMeta(QuerieSave)).
				WithArgs(tt.proj.OwnerId, tt.proj.Name).
				WillReturnError(tt.returnErr).
				WillReturnResult(sqlmock.NewResult(tt.returnLastId, tt.returnRowsAffected))

			postgres := NewPostgres(db)
			err = postgres.Save(context.Background(), tt.proj)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func TestPostgres_Delete(t *testing.T) {
	tests := []struct {
		testName string

		proj *projectdomain.ProjectDomain

		ownerId     uint32
		name        string
		rowAffected int64
		returnErr   error

		expErr error
	}{
		{
			testName: "Success",

			proj: &projectdomain.ProjectDomain{OwnerId: 1, Name: "Name"},

			ownerId:     1,
			name:        "Name",
			rowAffected: 1,
			returnErr:   nil,

			expErr: nil,
		}, {
			testName: "Not found",

			proj: &projectdomain.ProjectDomain{OwnerId: 1, Name: "Name"},

			ownerId:     1,
			name:        "Name",
			rowAffected: 0,
			returnErr:   nil,

			expErr: storage.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			mock.ExpectExec(regexp.QuoteMeta(QuerieDelete)).
				WithArgs(tt.ownerId, tt.name).
				WillReturnResult(sqlmock.NewResult(1, tt.rowAffected)).
				WillReturnError(tt.returnErr)

			postgres := NewPostgres(db)
			err = postgres.Delete(context.Background(), tt.proj)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
