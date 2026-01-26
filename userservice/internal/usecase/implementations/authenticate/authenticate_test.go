package authenticate

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"userservice/internal/repository/session"
	autherr "userservice/internal/usecase/errors/authenticate"
	authmocks "userservice/internal/usecase/implementations/authenticate/mocks"
	authmodel "userservice/internal/usecase/models/authenticate"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../repository/session/sessionrepo.go -destination=./mocks/mock_session.go -package=authmocks
func TestAuthenticate(t *testing.T) {
	tests := []struct {
		testName string

		sessionInput  string
		sessionOutput uint32
		sessionErr    error

		authInput *authmodel.AuthInput
		expOutput *authmodel.AuthOutput
		expErr    error
	}{
		{
			testName: "Success",

			sessionInput:  "sessionId",
			sessionOutput: 1,
			sessionErr:    nil,

			authInput: authmodel.NewAuthInput("sessionId"),
			expOutput: authmodel.NewAuthOutput(1),
			expErr:    nil,
		}, {
			testName: "Session not found",

			sessionInput:  "sessionId",
			sessionOutput: 0,
			sessionErr:    session.ErrKeyNotFound,

			authInput: authmodel.NewAuthInput("sessionId"),
			expOutput: authmodel.NewAuthOutput(0),
			expErr:    autherr.ErrSessionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			log := slog.New(slog.NewTextHandler(io.Discard, nil))
			sessionMock := authmocks.NewMockSessionRepo(ctrl)

			sessionMock.EXPECT().Get(gomock.Any(), tt.sessionInput).
				Return(tt.sessionOutput, tt.sessionErr)

			auth := NewGetUserIDBySessionUC(log, sessionMock)

			out, err := auth.Execute(context.Background(), tt.authInput)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expOutput, out)
		})
	}
}
