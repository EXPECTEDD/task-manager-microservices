package getownerid

import (
	"context"
	"io"
	"log/slog"
	projectdomain "projectservice/internal/domain/project"
	"projectservice/internal/repository/storage"
	getowneriderr "projectservice/internal/usecase/error/getownerid"
	getowneridmocks "projectservice/internal/usecase/implementations/getownerid/mocks"
	getowneridmodel "projectservice/internal/usecase/models/getownerid"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../repository/storage/storagerepo.go -destination=./mocks/mock_storage.go -package=getowneridmocks

func TestGetOwnerId(t *testing.T) {
	timeNow := time.Now()

	tests := []struct {
		testName string

		projectId        uint32
		storageReturn    *projectdomain.ProjectDomain
		storageReturnErr error

		input *getowneridmodel.GetOwnerIdInput

		expErr     error
		expOwnerId uint32
	}{
		{
			testName: "Success",

			projectId:        1,
			storageReturn:    projectdomain.RestoreProjectDomain(1, 1, "Proj", timeNow),
			storageReturnErr: nil,

			input: getowneridmodel.NewGetOwnerIdInput(1),

			expErr:     nil,
			expOwnerId: 1,
		}, {
			testName: "Not found",

			projectId:        1,
			storageReturn:    nil,
			storageReturnErr: storage.ErrNotFound,

			input: getowneridmodel.NewGetOwnerIdInput(1),

			expErr:     getowneriderr.ErrProjectsNotFound,
			expOwnerId: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storMock := getowneridmocks.NewMockStorageRepo(ctrl)
			storMock.EXPECT().GetProject(gomock.Any(), tt.projectId).
				Return(tt.storageReturn, tt.storageReturnErr)

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			getOwnerIdUC := NewGetOwnerIdUC(log, storMock)

			out, err := getOwnerIdUC.Execute(context.Background(), tt.input)
			require.Equal(t, tt.expErr, err)
			require.Equal(t, tt.expOwnerId, out.OwnerId)
		})
	}
}
