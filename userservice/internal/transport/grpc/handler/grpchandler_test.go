package grpchandler

import (
	"context"
	"io"
	"log/slog"
	"testing"
	grpchandlmocks "userservice/internal/transport/grpc/handler/mocks"
	autherr "userservice/internal/usecase/errors/authenticate"
	authmodel "userservice/internal/usecase/models/authenticate"
	userservicev1 "userservice/proto/userservice"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=./../../../usecase/interfaces/authenticate.go -destination=mocks/mock_authenticate.go -package=grpchandlmocks

func TestGRPCHandler(t *testing.T) {
	tests := []struct {
		testName string

		handlReq *userservicev1.GetIdBySessionRequest

		authInput  *authmodel.AuthInput
		authOutput *authmodel.AuthOutput
		authErr    error

		expOutput *userservicev1.GetIdBySessionResponse
		expErr    error
	}{
		{
			testName: "Success",

			handlReq: &userservicev1.GetIdBySessionRequest{
				SessionId: "sessionId",
			},

			authInput:  authmodel.NewAuthInput("sessionId"),
			authOutput: authmodel.NewAuthOutput(1),
			authErr:    nil,

			expOutput: &userservicev1.GetIdBySessionResponse{
				UserId: 1,
			},
			expErr: nil,
		}, {
			testName: "Session not found",

			handlReq: &userservicev1.GetIdBySessionRequest{
				SessionId: "sessionId",
			},

			authInput:  authmodel.NewAuthInput("sessionId"),
			authOutput: authmodel.NewAuthOutput(0),
			authErr:    autherr.ErrSessionNotFound,

			expOutput: nil,
			expErr:    status.Error(codes.NotFound, "session not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authUCMock := grpchandlmocks.NewMockGetUserIDBySessionUsecase(ctrl)

			authUCMock.EXPECT().Execute(gomock.Any(), tt.authInput).
				Return(tt.authOutput, tt.authErr)

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			grpcHandl := NewGRPCHandler(log, authUCMock)
			res, err := grpcHandl.GetIdBySession(context.Background(), tt.handlReq)
			assert.ErrorIs(t, err, tt.expErr)
			assert.Equal(t, tt.expOutput, res)
		})
	}
}
