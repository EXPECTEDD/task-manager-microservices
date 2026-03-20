package updateuc

import (
	"context"
	"io"
	"log/slog"
	"taskservice/internal/repository/storage"
	updatetaskerr "taskservice/internal/usecase/error/updatetask"
	updatemocks "taskservice/internal/usecase/implementations/updatetask/mocks"
	updatemodel "taskservice/internal/usecase/models/updatetask"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../repository/storage/storagerepo.go -destination=./mocks/mock_storage.go -package=updatemocks
func TestUpdateTask(t *testing.T) {
	timeNow := time.Now()

	tests := []struct {
		testName string

		expChangeDescription    bool
		newDescription          string
		changeDescriptionReturn error

		expChangeDeadline    bool
		newDeadline          time.Time
		changeDeadlineReturn error

		taskId    uint32
		projectId uint32

		in *updatemodel.UpdateTaskInput

		expReturn *updatemodel.UpdateTaskOutput
		expErr    error
	}{
		{
			testName: "Success",

			expChangeDescription:    true,
			newDescription:          "new description",
			changeDescriptionReturn: nil,

			expChangeDeadline:    true,
			newDeadline:          timeNow,
			changeDeadlineReturn: nil,

			taskId:    1,
			projectId: 1,

			expReturn: updatemodel.NewUpdateTaskOutput(true),
			expErr:    nil,
		}, {
			testName: "Success without description",

			expChangeDescription: false,

			expChangeDeadline:    true,
			newDeadline:          timeNow,
			changeDeadlineReturn: nil,

			taskId:    1,
			projectId: 1,

			expReturn: updatemodel.NewUpdateTaskOutput(true),
			expErr:    nil,
		}, {
			testName: "Success without deadline",

			expChangeDescription:    true,
			newDescription:          "new description",
			changeDescriptionReturn: nil,

			expChangeDeadline: false,

			taskId:    1,
			projectId: 1,

			expReturn: updatemodel.NewUpdateTaskOutput(true),
			expErr:    nil,
		}, {
			testName: "Success without all",

			expChangeDescription: false,

			expChangeDeadline: false,

			taskId:    1,
			projectId: 1,

			expReturn: updatemodel.NewUpdateTaskOutput(true),
			expErr:    nil,
		}, {
			testName: "Task not found",

			expChangeDescription:    true,
			newDescription:          "new description",
			changeDescriptionReturn: storage.ErrTaskNotFound,

			expChangeDeadline: false,

			taskId:    1,
			projectId: 1,

			expReturn: updatemodel.NewUpdateTaskOutput(false),
			expErr:    updatetaskerr.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt.in = &updatemodel.UpdateTaskInput{}
			tt.in.TaskId = tt.taskId
			tt.in.ProjectId = tt.projectId

			storMock := updatemocks.NewMockStorageRepo(ctrl)
			if tt.expChangeDescription {
				storMock.EXPECT().ChangeDescription(gomock.Any(), tt.taskId, tt.projectId, tt.newDescription).
					Return(tt.changeDescriptionReturn)
				tt.in.NewDescription = &tt.newDescription
			}
			if tt.expChangeDeadline {
				storMock.EXPECT().ChangeDeadline(gomock.Any(), tt.taskId, tt.projectId, tt.newDeadline).
					Return(tt.changeDeadlineReturn)
				tt.in.NewDeadline = &tt.newDeadline
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			updateUC := NewUpdateTaskUC(log, storMock)

			out, err := updateUC.Execute(context.Background(), tt.in)

			require.Equal(t, tt.expReturn.Updated, out.Updated)
			require.Equal(t, tt.expErr, err)
		})
	}
}
