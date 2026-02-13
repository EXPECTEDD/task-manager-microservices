package updateproject

import (
	"context"
	"io"
	"log/slog"
	"projectservice/internal/repository/storage"
	updateerr "projectservice/internal/usecase/error/updateproject"
	updatemocks "projectservice/internal/usecase/implementations/updateproject/mocks"
	updatemodel "projectservice/internal/usecase/models/updateproject"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../repository/storage/storagerepo.go -destination=./mocks/mock_storage.go -package=updatemocks
func TestUpdateProject(t *testing.T) {
	tests := []struct {
		testName string

		expStor       bool
		ownerId       uint32
		projectId     uint32
		newName       string
		storReturnErr error

		expErr error
		expOut *updatemodel.UpdateProjectOutput
	}{
		{
			testName: "Success",

			expStor:       true,
			ownerId:       1,
			projectId:     1,
			newName:       "newName",
			storReturnErr: nil,

			expErr: nil,
			expOut: updatemodel.NewUpdateProjectOutput(true),
		}, {
			testName: "Invalid new name",

			expStor:       false,
			ownerId:       1,
			projectId:     1,
			newName:       strings.Repeat("a", 256),
			storReturnErr: nil,

			expErr: updateerr.ErrInvalidName,
			expOut: updatemodel.NewUpdateProjectOutput(false),
		}, {
			testName: "Not found",

			expStor:       true,
			ownerId:       1,
			projectId:     1,
			newName:       "newName",
			storReturnErr: storage.ErrNotFound,

			expErr: updateerr.ErrProjectNotFound,
			expOut: updatemodel.NewUpdateProjectOutput(false),
		}, {
			testName: "Already exists",

			expStor:       true,
			ownerId:       1,
			projectId:     1,
			newName:       "newName",
			storReturnErr: storage.ErrAlreadyExists,

			expErr: updateerr.ErrProjectNameAlreadyExists,
			expOut: updatemodel.NewUpdateProjectOutput(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storMock := updatemocks.NewMockStorageRepo(ctrl)
			if tt.expStor {
				storMock.EXPECT().UpdateName(gomock.Any(), tt.ownerId, tt.projectId, tt.newName).
					Return(tt.storReturnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			updateUC := NewUpdateProjectUC(log, storMock)

			in := updatemodel.NewUpdateProjectInput(tt.ownerId, tt.projectId, &tt.newName)

			out, err := updateUC.Execute(context.Background(), in)
			require.Equal(t, tt.expErr, err)
			require.Equal(t, tt.expOut, out)
		})
	}
}
