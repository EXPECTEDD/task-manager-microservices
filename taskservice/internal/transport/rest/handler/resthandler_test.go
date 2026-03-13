package resthandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	handlmocks "taskservice/internal/transport/rest/handler/mocks"
	createmodel "taskservice/internal/usecase/models/createtask"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../usecase/interfaces/create_task.go -destination=./mocks/mock_create_project.go -package=handlmocks
func TestResthandler_Create(t *testing.T) {
	timeNow := time.Now().Round(0)

	tests := []struct {
		testName string

		projectId uint32

		expCreateMock   bool
		createIn        *createmodel.CreateTaskInput
		createReturnOut *createmodel.CreateTaskOutput
		createReturnErr error

		body map[string]any

		expTaskId     uint32
		expStatusCode int
	}{
		{
			testName: "Success",

			projectId: 1,

			expCreateMock: true,
			createIn: createmodel.NewCreateInput(
				1,
				"desc",
				timeNow,
			),
			createReturnOut: createmodel.NewCreateOutput(
				1,
			),
			createReturnErr: nil,

			body: map[string]any{
				"description": "desc",
				"deadline":    timeNow,
			},

			expTaskId:     1,
			expStatusCode: http.StatusOK,
		}, {
			testName: "Invalid project id",

			projectId: 0,

			expCreateMock:   false,
			createIn:        nil,
			createReturnOut: nil,
			createReturnErr: nil,

			body: map[string]any{
				"description": "desc",
				"deadline":    timeNow,
			},

			expTaskId:     0,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Missing filed description",

			projectId: 1,

			expCreateMock:   false,
			createIn:        nil,
			createReturnOut: nil,
			createReturnErr: nil,

			body: map[string]any{
				"desc":     "desc",
				"deadline": timeNow,
			},

			expTaskId:     0,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Missing filed deadline",

			projectId: 1,

			expCreateMock: true,
			createIn: createmodel.NewCreateInput(
				1,
				"desc",
				time.Time{}.Round(0),
			),
			createReturnOut: createmodel.NewCreateOutput(
				1,
			),
			createReturnErr: nil,

			body: map[string]any{
				"description": "desc",
			},

			expTaskId:     1,
			expStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			createUCmock := handlmocks.NewMockCreateTaskUsecase(ctrl)
			if tt.expCreateMock {
				createUCmock.EXPECT().Execute(gomock.Any(), tt.createIn).
					Return(tt.createReturnOut, tt.createReturnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			handl := NewRestHandler(log, createUCmock)

			router := gin.New()
			router.POST("/test/:project_id", handl.Create)

			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.body)

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/test/%d", tt.projectId), bytes.NewReader(b))
			require.NoError(t, err)
			defer req.Body.Close()

			router.ServeHTTP(w, req)

			var respBody struct {
				TaskId uint32 `json:"task_id"`
			}

			require.NoError(t, json.NewDecoder(w.Body).Decode(&respBody))
			require.Equal(t, tt.expTaskId, respBody.TaskId)
			require.Equal(t, tt.expStatusCode, w.Result().StatusCode)
		})
	}
}
