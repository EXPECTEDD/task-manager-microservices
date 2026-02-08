package postgres

import (
	"context"
	projectdomain "projectservice/internal/domain/project"
	"projectservice/internal/repository/storage"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgres_Save(t *testing.T) {
	tests := []struct {
		testName string

		proj *projectdomain.ProjectDomain

		returnRows sqlmock.Rows
		returnErr  error
		expErr     error
		expId      uint32
	}{
		{
			testName: "Success",

			proj: &projectdomain.ProjectDomain{OwnerId: 1, Name: "Name"},

			returnRows: *sqlmock.NewRows([]string{"id"}).AddRow(1),
			returnErr:  nil,
			expErr:     nil,
			expId:      1,
		}, {
			testName: "Already exists",

			proj: &projectdomain.ProjectDomain{OwnerId: 1, Name: "Name"},

			returnRows: *sqlmock.NewRows([]string{"id"}).AddRow(0),
			returnErr:  &pq.Error{Code: "23505"},
			expErr:     storage.ErrAlreadyExists,
			expId:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta(QuerieSave)).
				WithArgs(tt.proj.OwnerId, tt.proj.Name).
				WillReturnError(tt.returnErr).
				WillReturnRows(&tt.returnRows)

			postgres := NewPostgres(db)
			id, err := postgres.Save(context.Background(), tt.proj)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expId, id)
		})
	}
}

func TestPostgres_Delete(t *testing.T) {
	tests := []struct {
		testName string

		projectId   uint32
		rowAffected int64
		returnErr   error

		expErr error
	}{
		{
			testName: "Success",

			projectId:   1,
			rowAffected: 1,
			returnErr:   nil,

			expErr: nil,
		}, {
			testName: "Not found",

			projectId:   1,
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
				WithArgs(tt.projectId).
				WillReturnResult(sqlmock.NewResult(1, tt.rowAffected)).
				WillReturnError(tt.returnErr)

			postgres := NewPostgres(db)
			err = postgres.Delete(context.Background(), tt.projectId)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func TestPostgres_GetAll(t *testing.T) {
	timeNow := time.Now()

	tests := []struct {
		testName   string
		ownerId    uint32
		returnRows *sqlmock.Rows
		returnErr  error
		expOutput  []*projectdomain.ProjectDomain
		expErr     error
	}{
		{
			testName: "Success",
			ownerId:  1,
			returnRows: sqlmock.NewRows([]string{
				"id", "owner_id", "name", "created_at",
			}).AddRow(1, 1, "A", timeNow),
			returnErr: nil,
			expOutput: []*projectdomain.ProjectDomain{
				{Id: 1, OwnerId: 1, Name: "A", CreatedAt: timeNow},
			},
			expErr: nil,
		}, {
			testName: "More returned projects",
			ownerId:  1,
			returnRows: sqlmock.NewRows([]string{
				"id", "owner_id", "name", "created_at",
			}).AddRow(1, 1, "A", timeNow).
				AddRow(2, 1, "B", timeNow).
				AddRow(3, 1, "C", timeNow).
				AddRow(4, 1, "D", timeNow).
				AddRow(5, 1, "E", timeNow).
				AddRow(6, 1, "F", timeNow).
				AddRow(7, 1, "G", timeNow).
				AddRow(8, 1, "H", timeNow).
				AddRow(9, 1, "I", timeNow).
				AddRow(10, 1, "J", timeNow),
			returnErr: nil,
			expOutput: []*projectdomain.ProjectDomain{
				{Id: 1, OwnerId: 1, Name: "A", CreatedAt: timeNow},
				{Id: 2, OwnerId: 1, Name: "B", CreatedAt: timeNow},
				{Id: 3, OwnerId: 1, Name: "C", CreatedAt: timeNow},
				{Id: 4, OwnerId: 1, Name: "D", CreatedAt: timeNow},
				{Id: 5, OwnerId: 1, Name: "E", CreatedAt: timeNow},
				{Id: 6, OwnerId: 1, Name: "F", CreatedAt: timeNow},
				{Id: 7, OwnerId: 1, Name: "G", CreatedAt: timeNow},
				{Id: 8, OwnerId: 1, Name: "H", CreatedAt: timeNow},
				{Id: 9, OwnerId: 1, Name: "I", CreatedAt: timeNow},
				{Id: 10, OwnerId: 1, Name: "J", CreatedAt: timeNow},
			},
			expErr: nil,
		}, {
			testName: "Not found",
			ownerId:  1,
			returnRows: sqlmock.NewRows([]string{
				"id", "owner_id", "name", "created_at",
			}),
			expOutput: nil,
			expErr:    storage.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta(QuerieGetAll)).
				WithArgs(tt.ownerId).
				WillReturnRows(tt.returnRows).
				WillReturnError(tt.returnErr)

			postgres := NewPostgres(db)

			projects, err := postgres.GetAll(context.Background(), tt.ownerId)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expOutput, projects)
		})
	}
}
