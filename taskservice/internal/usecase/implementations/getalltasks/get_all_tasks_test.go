package getalluc

import (
	"context"
	"io"
	"log/slog"
	taskdomain "taskservice/internal/domain/task"
	"taskservice/internal/repository/storage"
	getallerr "taskservice/internal/usecase/error/getalltasks"
	getallmocks "taskservice/internal/usecase/implementations/getalltasks/mocks"
	getallmodel "taskservice/internal/usecase/models/getalltasks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../repository/storage/storagerepo.go -destination=./mocks/mock_storage.go -package=getallmocks
func TestGetAllTasks(t *testing.T) {
	timeNow := time.Now().Round(1)

	tests := []struct {
		testName     string
		projectId    uint32
		expReturn    []*taskdomain.TaskDomain
		expReturnErr error

		expOut *getallmodel.GetAllTasksOutput
		expErr error
	}{
		{
			testName:     "Success",
			projectId:    1,
			expReturn:    []*taskdomain.TaskDomain{{Id: 1, ProjectId: 1, Description: "asd", Deadline: timeNow}, {Id: 2, ProjectId: 1, Description: "dsa", Deadline: timeNow}},
			expReturnErr: nil,

			expOut: &getallmodel.GetAllTasksOutput{Tasks: []*taskdomain.TaskDomain{{Id: 1, ProjectId: 1, Description: "asd", Deadline: timeNow}, {Id: 2, ProjectId: 1, Description: "dsa", Deadline: timeNow}}},
			expErr: nil,
		}, {
			testName:     "Tasks not found",
			projectId:    1,
			expReturn:    []*taskdomain.TaskDomain{},
			expReturnErr: storage.ErrTasksNotFound,

			expOut: &getallmodel.GetAllTasksOutput{Tasks: nil},
			expErr: getallerr.ErrTasksNotFound,
		}, {
			testName:     "Invalid project id",
			projectId:    0,
			expReturn:    []*taskdomain.TaskDomain{},
			expReturnErr: storage.ErrTasksNotFound,

			expOut: &getallmodel.GetAllTasksOutput{Tasks: nil},
			expErr: getallerr.ErrTasksNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storMock := getallmocks.NewMockStorageRepo(ctrl)
			storMock.EXPECT().GetAll(gomock.Any(), tt.projectId).
				Return(tt.expReturn, tt.expReturnErr)

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			getAllUC := NewGetAllTasksUC(log, storMock)

			out, err := getAllUC.Execute(context.Background(), getallmodel.NewGetAllTasksInput(tt.projectId))
			require.Equal(t, tt.expOut, out)
			require.Equal(t, tt.expErr, err)
		})
	}
}
