package deletetask

import (
	"context"
	"io"
	"log/slog"
	"taskservice/internal/repository/storage"
	deleteerr "taskservice/internal/usecase/error/deletetask"
	deletemocks "taskservice/internal/usecase/implementations/deletetask/mocks"
	deletemodel "taskservice/internal/usecase/models/deletetask"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../repository/storage/storagerepo.go -destination=./mocks/mock_storage.go -package=deletemocks
func TestDelete(t *testing.T) {
	tests := []struct {
		testName string

		taskId       uint32
		projectId    uint32
		expReturnErr error

		expOut bool
		expErr error
	}{
		{
			testName: "Success",

			taskId:       1,
			projectId:    1,
			expReturnErr: nil,

			expOut: true,
			expErr: nil,
		}, {
			testName: "Task not found",

			taskId:       1,
			projectId:    1,
			expReturnErr: storage.ErrTaskNotFound,

			expOut: false,
			expErr: deleteerr.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storMock := deletemocks.NewMockStorageRepo(ctrl)
			storMock.EXPECT().Delete(gomock.Any(), tt.taskId, tt.projectId).
				Return(tt.expReturnErr)

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			deleteUC := NewDeleteUC(log, storMock)

			in := deletemodel.NewDeleteTaskInput(
				tt.taskId,
				tt.projectId,
			)

			out, err := deleteUC.Execute(context.Background(), in)
			require.Equal(t, tt.expOut, out.Deleted)
			require.Equal(t, tt.expErr, err)
		})
	}
}
