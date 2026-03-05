package grpchandler

import (
	"context"
	"io"
	"log/slog"
	grpchandlermocks "projectservice/internal/transport/grpc/handler/mocks"
	getowneriderr "projectservice/internal/usecase/error/getownerid"
	getowneridmodel "projectservice/internal/usecase/models/getownerid"
	projectservicev1 "projectservice/proto/projectservice"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=./../../../usecase/interfaces/get_owner_id.go -destination=./mocks/mock_getownerid.go -package=grpchandlermocks
func TestGrpcHandler_GetOwnerId(t *testing.T) {
	tests := []struct {
		testName string

		expGetOwnerId       bool
		input               *getowneridmodel.GetOwnerIdInput
		getOwnerIdReturn    *getowneridmodel.GetOwnerIdOutput
		getOwnerIdReturnErr error

		req        *projectservicev1.GetOwnerIdRequest
		expErr     error
		expOwnerId *projectservicev1.GetOwnerIdResponse
	}{
		{
			testName: "Success",

			expGetOwnerId:       true,
			input:               getowneridmodel.NewGetOwnerIdInput(1),
			getOwnerIdReturn:    getowneridmodel.NewGetOwnerIdOutput(1),
			getOwnerIdReturnErr: nil,

			req: &projectservicev1.GetOwnerIdRequest{
				ProjectId: 1,
			},
			expErr: nil,
			expOwnerId: &projectservicev1.GetOwnerIdResponse{
				OwnerId: 1,
			},
		}, {
			testName: "Project not found",

			expGetOwnerId:       true,
			input:               getowneridmodel.NewGetOwnerIdInput(1),
			getOwnerIdReturn:    getowneridmodel.NewGetOwnerIdOutput(0),
			getOwnerIdReturnErr: getowneriderr.ErrProjectsNotFound,

			req: &projectservicev1.GetOwnerIdRequest{
				ProjectId: 1,
			},
			expErr:     status.Error(codes.NotFound, "project not found"),
			expOwnerId: nil,
		}, {
			testName: "Ivalid project id",

			expGetOwnerId: false,

			req: &projectservicev1.GetOwnerIdRequest{
				ProjectId: 0,
			},
			expErr:     status.Error(codes.InvalidArgument, "invalid project id"),
			expOwnerId: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			getOwnerIdUCMock := grpchandlermocks.NewMockGetOwnerIdUsecase(ctrl)
			if tt.expGetOwnerId {
				getOwnerIdUCMock.EXPECT().Execute(gomock.Any(), tt.input).
					Return(tt.getOwnerIdReturn, tt.getOwnerIdReturnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			handler := NewGRPCServer(log, getOwnerIdUCMock)

			res, err := handler.GetOwnerId(context.Background(), tt.req)
			require.Equal(t, tt.expErr, err)
			require.Equal(t, tt.expOwnerId, res)
		})
	}
}
