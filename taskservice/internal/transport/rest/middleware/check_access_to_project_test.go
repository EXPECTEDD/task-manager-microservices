package middleware

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	middlewaremocks "taskservice/internal/transport/rest/middleware/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=./../../../repository/projectrepository/project_repository.go -destination=./mocks/mock_project_repository.go -package=middlewaremocks
func TestCheckAccessToProjectMiddleware(t *testing.T) {
	tests := []struct {
		testName string

		sessionId string
		userId    uint32

		expProjRep       bool
		projectId        uint32
		projRepReturn    uint32
		projRepReturnErr error

		expStatusCode int
	}{
		{
			testName: "Success",

			sessionId: "session",
			userId:    1,
			projectId: 1,

			expProjRep:       true,
			projRepReturn:    1,
			projRepReturnErr: nil,

			expStatusCode: http.StatusOK,
		}, {
			testName: "Invalid project id",

			sessionId: "session",
			userId:    1,
			projectId: 0,

			expProjRep: false,

			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Project repository internal error",

			sessionId: "session",
			userId:    1,
			projectId: 1,

			expProjRep:       true,
			projRepReturn:    0,
			projRepReturnErr: status.Error(codes.Internal, "error"),

			expStatusCode: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sessionValidMock := middlewaremocks.NewMockSessionValidator(ctrl)
			sessionValidMock.EXPECT().GetIdBySession(gomock.Any(), tt.sessionId).
				Return(tt.userId, nil)

			projRepMock := middlewaremocks.NewMockProjectRepository(ctrl)
			if tt.expProjRep {
				projRepMock.EXPECT().GetOwnerId(gomock.Any(), tt.projectId).
					Return(tt.projRepReturn, tt.projRepReturnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			router := gin.New()
			router.Use(GetSessionMiddleware(log))
			router.Use(SessionAuthMiddleware(log, sessionValidMock, 1*time.Second))
			router.Use(CheckAccessToProjectMiddleware(log, projRepMock, 2*time.Millisecond))
			router.GET("/test/:project_id", func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/test/%d", tt.projectId), nil)

			req.AddCookie(&http.Cookie{
				Name:  "sessionId",
				Value: tt.sessionId,
			})

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, tt.expStatusCode, w.Result().StatusCode)
		})
	}
}
