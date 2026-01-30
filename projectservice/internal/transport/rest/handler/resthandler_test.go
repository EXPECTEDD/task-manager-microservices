package resthandler

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	projectdomain "projectservice/internal/domain/project"
	resthandlmocks "projectservice/internal/transport/rest/handler/mocks"
	"projectservice/internal/transport/rest/middleware"
	createerr "projectservice/internal/usecase/error/createproject"
	createmodel "projectservice/internal/usecase/models/createproject"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../repository/sessionvalidator/session_validator.go -destination=./mocks/mock_session_validator.go -package=resthandlmocks
//go:generate mockgen -source=./../../../usecase/interfaces/create_project.go -destination=./mocks/mock_create_project.go -package=resthandlmocks
func TestRestHandler_Create(t *testing.T) {
	tests := []struct {
		testName string

		sessionId string
		userId    uint32

		expCreate       bool
		createInput     *createmodel.CreateProjectInput
		createOutput    *createmodel.CreateProjectOutput
		createReturnErr error

		contentType string
		body        map[string]string

		expResp       bool
		expStatusCode int
	}{
		{
			testName: "Success",

			sessionId: "sessionId",
			userId:    1,

			expCreate:       true,
			createInput:     createmodel.NewCreateProjectInput(1, "Name"),
			createOutput:    createmodel.NewCreateProjectOutput(true),
			createReturnErr: nil,

			contentType: "application/json",
			body: map[string]string{
				"name": "Name",
			},

			expResp:       true,
			expStatusCode: http.StatusOK,
		}, {
			testName: "Missing field name",

			sessionId: "sessionId",
			userId:    1,

			expCreate:       false,
			createInput:     createmodel.NewCreateProjectInput(1, "Name"),
			createOutput:    createmodel.NewCreateProjectOutput(true),
			createReturnErr: nil,

			contentType: "application/json",
			body: map[string]string{
				"nme": "Name",
			},

			expResp:       false,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Invalid name",

			sessionId: "sessionId",
			userId:    1,

			expCreate:       true,
			createInput:     createmodel.NewCreateProjectInput(1, strings.Repeat("Name", 300)),
			createOutput:    createmodel.NewCreateProjectOutput(false),
			createReturnErr: projectdomain.ErrInvalidName,

			contentType: "application/json",
			body: map[string]string{
				"name": strings.Repeat("Name", 300),
			},

			expResp:       false,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Already exists",

			sessionId: "sessionId",
			userId:    1,

			expCreate:       true,
			createInput:     createmodel.NewCreateProjectInput(1, "Name"),
			createOutput:    createmodel.NewCreateProjectOutput(false),
			createReturnErr: createerr.ErrAlreadyExists,

			contentType: "application/json",
			body: map[string]string{
				"name": "Name",
			},

			expResp:       false,
			expStatusCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			createUC := resthandlmocks.NewMockCreateProjectUsecase(ctrl)
			if tt.expCreate {
				createUC.EXPECT().Execute(gomock.Any(), tt.createInput).
					Return(tt.createOutput, tt.createReturnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			client := resthandlmocks.NewMockSessionValidator(ctrl)
			client.EXPECT().GetIdBySession(gomock.Any(), tt.sessionId).
				Return(tt.userId, nil)

			handl := NewHandler(log, createUC)

			router := gin.New()
			router.Use(middleware.GetSessionMiddleware(log))
			router.Use(middleware.SessionAuthMiddleware(log, client, 10*time.Second))
			router.POST("/test", handl.Create)

			b, err := json.Marshal(tt.body)

			req, err := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(b))
			assert.NoError(t, err)

			c := &http.Cookie{
				Name:  "sessionId",
				Value: tt.sessionId,
			}
			req.AddCookie(c)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			var respBody struct {
				IsCreated bool `json:"is_created"`
			}

			assert.NoError(t, json.NewDecoder(w.Body).Decode(&respBody))
			assert.Equal(t, tt.expResp, respBody.IsCreated)
			assert.Equal(t, tt.expStatusCode, w.Result().StatusCode)
		})
	}
}
