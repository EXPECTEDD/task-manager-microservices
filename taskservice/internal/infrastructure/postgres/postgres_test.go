package postgres

import (
	"context"
	"database/sql/driver"
	"regexp"
	taskdomain "taskservice/internal/domain/task"
	posmodels "taskservice/internal/infrastructure/postgres/models"
	"taskservice/internal/repository/storage"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestPostgres_Save_Success(t *testing.T) {
	timeNow := time.Now()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	posModel := posmodels.NewTaskPosModel(
		0,
		1,
		"desc",
		timeNow,
	)

	mock.ExpectQuery(regexp.QuoteMeta(QuerieCreate)).
		WithArgs(posModel.ProjectId, posModel.Description, posModel.Deadline).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)).
		WillReturnError(nil)

	postgres := NewPostgres(db)

	td, err := taskdomain.NewTaskDomain(
		1,
		"desc",
		timeNow,
	)
	require.NoError(t, err)

	id, err := postgres.Save(context.Background(), td)
	require.NoError(t, err)
	require.Equal(t, uint32(1), id)
}

func TestPostgres_ChangeDescription(t *testing.T) {
	tests := []struct {
		testName string

		taskId         uint32
		projectId      uint32
		newDescription string
		returnResult   driver.Result

		expectErr error
	}{
		{
			testName: "Success",

			taskId:         1,
			projectId:      1,
			newDescription: "new description",
			returnResult:   sqlmock.NewResult(0, 1),

			expectErr: nil,
		}, {
			testName: "Task not found",

			taskId:         1,
			projectId:      1,
			newDescription: "new description",
			returnResult:   sqlmock.NewResult(0, 0),

			expectErr: storage.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectExec(regexp.QuoteMeta(QuerieUpdateDescription)).
				WithArgs(tt.newDescription, tt.taskId, tt.projectId).
				WillReturnResult(tt.returnResult).
				WillReturnError(nil)

			postgres := NewPostgres(db)

			err = postgres.ChangeDescription(context.Background(), tt.taskId, tt.projectId, tt.newDescription)

			require.Equal(t, tt.expectErr, err)
		})
	}
}

func TestPostgres_ChangeDeadline(t *testing.T) {
	timeNow := time.Now()

	tests := []struct {
		testName string

		taskId       uint32
		projectId    uint32
		newDeadline  time.Time
		returnResult driver.Result

		expectErr error
	}{
		{
			testName: "Success",

			taskId:       1,
			projectId:    1,
			newDeadline:  timeNow,
			returnResult: sqlmock.NewResult(0, 1),

			expectErr: nil,
		}, {
			testName: "Task not found",

			taskId:       1,
			projectId:    1,
			newDeadline:  timeNow,
			returnResult: sqlmock.NewResult(0, 0),

			expectErr: storage.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectExec(regexp.QuoteMeta(QuerieUpdateDeadline)).
				WithArgs(tt.newDeadline, tt.taskId, tt.projectId).
				WillReturnResult(tt.returnResult).
				WillReturnError(nil)

			postgres := NewPostgres(db)

			err = postgres.ChangeDeadline(context.Background(), tt.taskId, tt.projectId, tt.newDeadline)

			require.Equal(t, tt.expectErr, err)
		})
	}
}

func TestPostgres_Delete(t *testing.T) {
	tests := []struct {
		testName string

		taskId       uint32
		projectId    uint32
		returnResult driver.Result

		expectErr error
	}{
		{
			testName: "Success",

			taskId:       1,
			projectId:    1,
			returnResult: sqlmock.NewResult(0, 1),

			expectErr: nil,
		}, {
			testName: "Task not found",

			taskId:       1,
			projectId:    1,
			returnResult: sqlmock.NewResult(0, 0),

			expectErr: storage.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectExec(regexp.QuoteMeta(QuerieDelete)).
				WithArgs(tt.taskId, tt.projectId).
				WillReturnResult(tt.returnResult).
				WillReturnError(nil)

			postgres := NewPostgres(db)

			err = postgres.Delete(context.Background(), tt.taskId, tt.projectId)

			require.Equal(t, tt.expectErr, err)
		})
	}
}

func TestPostgres_GetAll(t *testing.T) {
	timeNow := time.Now().Round(1)

	tests := []struct {
		testName   string
		projectId  uint32
		returnRows *sqlmock.Rows
		returnErr  error

		expTasks []*taskdomain.TaskDomain
		expErr   error
	}{
		{
			testName:   "Success",
			projectId:  1,
			returnRows: sqlmock.NewRows([]string{"id", "project_id", "description", "deadline"}).AddRow(1, 1, "asd", timeNow).AddRow(2, 1, "dsa", timeNow),
			returnErr:  nil,

			expTasks: []*taskdomain.TaskDomain{{Id: 1, ProjectId: 1, Description: "asd", Deadline: timeNow}, {Id: 2, ProjectId: 1, Description: "dsa", Deadline: timeNow}},
			expErr:   nil,
		}, {
			testName:   "Tasks not found",
			projectId:  1,
			returnRows: sqlmock.NewRows([]string{"id", "project_id", "description", "deadline"}),
			returnErr:  nil,

			expTasks: nil,
			expErr:   storage.ErrTasksNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta(QuerieGetAll)).
				WithArgs(tt.projectId).
				WillReturnRows(tt.returnRows).
				WillReturnError(tt.returnErr)

			postgres := NewPostgres(db)

			tasks, err := postgres.GetAll(context.Background(), tt.projectId)
			require.Equal(t, tt.expTasks, tasks)
			require.Equal(t, tt.expErr, err)
		})
	}
}

func TestPostgres_Get(t *testing.T) {
	timeNow := time.Now().Round(1)

	tests := []struct {
		testName   string
		taskId     uint32
		projectId  uint32
		returnRows *sqlmock.Rows
		returnErr  error

		expTasks *taskdomain.TaskDomain
		expErr   error
	}{
		{
			testName:   "Success",
			taskId:     1,
			projectId:  1,
			returnRows: sqlmock.NewRows([]string{"id", "project_id", "description", "deadline"}).AddRow(1, 1, "asd", timeNow),
			returnErr:  nil,

			expTasks: &taskdomain.TaskDomain{Id: 1, ProjectId: 1, Description: "asd", Deadline: timeNow},
			expErr:   nil,
		}, {
			testName:   "Tasks not found",
			taskId:     1,
			projectId:  1,
			returnRows: sqlmock.NewRows([]string{"id", "project_id", "description", "deadline"}),
			returnErr:  nil,

			expTasks: nil,
			expErr:   storage.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta(QuerieGet)).
				WithArgs(tt.projectId, tt.taskId).
				WillReturnRows(tt.returnRows).
				WillReturnError(tt.returnErr)

			postgres := NewPostgres(db)

			task, err := postgres.Get(context.Background(), tt.projectId, tt.taskId)
			require.Equal(t, tt.expTasks, task)
			require.Equal(t, tt.expErr, err)
		})
	}
}
