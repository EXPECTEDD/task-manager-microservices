package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	middlewaremocks "projectservice/internal/transport/rest/middleware/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=./../../../repository/sessionvalidator/session_validator.go -destination=./mocks/mock_session_validator.go -package=middlewaremocks

func TestSessionAuthMiddleware(t *testing.T) {
	tests := []struct {
		testName string

		sessIdInput string
		sessOutput  uint32
		sessErr     error

		expCode int
	}{
		{
			testName: "Success",

			sessIdInput: "123321",
			sessOutput:  1,
			sessErr:     nil,

			expCode: http.StatusOK,
		}, {
			testName: "Session not found",

			sessIdInput: "123321",
			sessOutput:  0,
			sessErr:     status.Error(codes.NotFound, "session not found"),

			expCode: http.StatusNotFound,
		}, {
			testName: "User service internal error",

			sessIdInput: "123321",
			sessOutput:  0,
			sessErr:     status.Error(codes.Internal, "internal server error"),

			expCode: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sessValidMock := middlewaremocks.NewMockSessionValidator(ctrl)
			sessValidMock.EXPECT().GetIdBySession(gomock.Any(), tt.sessIdInput).
				Return(tt.sessOutput, tt.sessErr)

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			gin.SetMode(gin.DebugMode)
			router := gin.New()
			router.Use(GetSessionMiddleware(log))
			router.Use(SessionAuthMiddleware(log, sessValidMock, 15*time.Second))
			router.GET("/test", func(ctx *gin.Context) {
				userId, ok := ctx.Get("userId")
				assert.True(t, ok)
				assert.Equal(t, tt.sessOutput, userId)
				ctx.JSON(http.StatusOK, "ok")
			})

			req, err := http.NewRequest(http.MethodGet, "/test", nil)
			assert.NoError(t, err)

			req.AddCookie(&http.Cookie{
				Name:  "sessionId",
				Value: tt.sessIdInput,
			})

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expCode, w.Result().StatusCode)
		})
	}
}
