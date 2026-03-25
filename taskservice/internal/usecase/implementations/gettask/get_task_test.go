package getuc

import (
	"context"
	"io"
	"log/slog"
	taskdomain "taskservice/internal/domain/task"
	"taskservice/internal/repository/storage"
	geterr "taskservice/internal/usecase/error/gettask"
	getmocks "taskservice/internal/usecase/implementations/gettask/mocks"
	getmodel "taskservice/internal/usecase/models/gettask"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../repository/storage/storagerepo.go -destination=./mocks/mock_storage.go -package=getmocks
func TestGetAllTasks(t *testing.T) {
	timeNow := time.Now().Round(1)

	tests := []struct {
		testName     string
		projectId    uint32
		taskId       uint32
		expReturn    *taskdomain.TaskDomain
		expReturnErr error

		expOut *getmodel.GetTaskOutput
		expErr error
	}{
		{
			testName:     "Success",
			projectId:    1,
			taskId:       1,
			expReturn:    &taskdomain.TaskDomain{Id: 1, ProjectId: 1, Description: "asd", Deadline: timeNow},
			expReturnErr: nil,

			expOut: &getmodel.GetTaskOutput{Task: &taskdomain.TaskDomain{Id: 1, ProjectId: 1, Description: "asd", Deadline: timeNow}},
			expErr: nil,
		}, {
			testName:     "Task not found",
			projectId:    1,
			taskId:       1,
			expReturn:    &taskdomain.TaskDomain{},
			expReturnErr: storage.ErrTaskNotFound,

			expOut: &getmodel.GetTaskOutput{Task: nil},
			expErr: geterr.ErrTaskNotFound,
		}, {
			testName:     "Invalid project id",
			projectId:    0,
			taskId:       1,
			expReturn:    &taskdomain.TaskDomain{},
			expReturnErr: storage.ErrTaskNotFound,

			expOut: &getmodel.GetTaskOutput{Task: nil},
			expErr: geterr.ErrTaskNotFound,
		}, {
			testName:     "Invalid task id",
			projectId:    1,
			taskId:       0,
			expReturn:    &taskdomain.TaskDomain{},
			expReturnErr: storage.ErrTaskNotFound,

			expOut: &getmodel.GetTaskOutput{Task: nil},
			expErr: geterr.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storMock := getmocks.NewMockStorageRepo(ctrl)
			storMock.EXPECT().Get(gomock.Any(), tt.taskId, tt.projectId).
				Return(tt.expReturn, tt.expReturnErr)

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			getUC := NewGetTaskUC(log, storMock)

			out, err := getUC.Execute(context.Background(), getmodel.NewGetTaskInput(tt.taskId, tt.projectId))
			require.Equal(t, tt.expOut, out)
			require.Equal(t, tt.expErr, err)
		})
	}
}
