package resthandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	taskdomain "taskservice/internal/domain/task"
	handlmocks "taskservice/internal/transport/rest/handler/mocks"
	deleteerr "taskservice/internal/usecase/error/deletetask"
	getallerr "taskservice/internal/usecase/error/getalltasks"
	geterr "taskservice/internal/usecase/error/gettask"
	createmodel "taskservice/internal/usecase/models/createtask"
	deletemodel "taskservice/internal/usecase/models/deletetask"
	getallmodel "taskservice/internal/usecase/models/getalltasks"
	getmodel "taskservice/internal/usecase/models/gettask"
	updatemodel "taskservice/internal/usecase/models/updatetask"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../usecase/interfaces/create_task.go -destination=./mocks/mock_create_task.go -package=handlmocks
//go:generate mockgen -source=./../../../usecase/interfaces/update_task.go -destination=./mocks/mock_update_task.go -package=handlmocks
//go:generate mockgen -source=./../../../usecase/interfaces/delete_task.go -destination=./mocks/mock_delete_task.go -package=handlmocks
//go:generate mockgen -source=./../../../usecase/interfaces/get_all_tasks.go -destination=./mocks/mock_get_all_tasks.go -package=handlmocks
//go:generate mockgen -source=./../../../usecase/interfaces/get_task.go	-destination=./mocks/mock_get_task.go -package=handlmocks
func TestRestHandler_Create(t *testing.T) {
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

			handl := NewRestHandler(log, createUCmock, nil, nil, nil, nil)

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

func TestRestHandler_Update(t *testing.T) {
	timeNow := time.Now()

	tests := []struct {
		testName string

		taskId    uint32
		projectId uint32

		body map[string]string

		expUpdateMock       bool
		updateMockReturn    *updatemodel.UpdateTaskOutput
		updateMockReturnErr error

		expReturn     bool
		expStatusCode int
	}{
		{
			testName: "Success",

			taskId:    1,
			projectId: 1,

			body: map[string]string{
				"new_description": "new description",
				"new_deadline":    timeNow.Format(time.RFC3339),
			},

			expUpdateMock:       true,
			updateMockReturn:    updatemodel.NewUpdateTaskOutput(true),
			updateMockReturnErr: nil,

			expReturn:     true,
			expStatusCode: http.StatusOK,
		}, {
			testName: "Succes without description",

			taskId:    1,
			projectId: 1,

			body: map[string]string{
				"new_deadline": timeNow.Format(time.RFC3339),
			},

			expUpdateMock:       true,
			updateMockReturn:    updatemodel.NewUpdateTaskOutput(true),
			updateMockReturnErr: nil,

			expReturn:     true,
			expStatusCode: http.StatusOK,
		}, {
			testName: "Success without deadline",

			taskId:    1,
			projectId: 1,

			body: map[string]string{
				"new_description": "new description",
			},

			expUpdateMock:       true,
			updateMockReturn:    updatemodel.NewUpdateTaskOutput(true),
			updateMockReturnErr: nil,

			expReturn:     true,
			expStatusCode: http.StatusOK,
		}, {
			testName: "Without all",

			taskId:    1,
			projectId: 1,

			body: map[string]string{},

			expUpdateMock:       false,
			updateMockReturn:    updatemodel.NewUpdateTaskOutput(true),
			updateMockReturnErr: nil,

			expReturn:     false,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Invalid task id",

			taskId:    0,
			projectId: 1,

			body: map[string]string{
				"new_description": "new description",
			},

			expUpdateMock:       false,
			updateMockReturn:    updatemodel.NewUpdateTaskOutput(true),
			updateMockReturnErr: nil,

			expReturn:     false,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Invalid project id",

			taskId:    1,
			projectId: 0,

			body: map[string]string{
				"new_description": "new description",
			},

			expUpdateMock:       false,
			updateMockReturn:    updatemodel.NewUpdateTaskOutput(true),
			updateMockReturnErr: nil,

			expReturn:     false,
			expStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			updateMock := handlmocks.NewMockUpdateTaskUsecase(ctrl)
			if tt.expUpdateMock {
				updateMock.EXPECT().Execute(gomock.Any(), gomock.Any()).
					Return(tt.updateMockReturn, tt.updateMockReturnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			handl := NewRestHandler(log, nil, updateMock, nil, nil, nil)

			router := gin.New()
			router.PATCH("/test/:task_id/:project_id", handl.Update)

			w := httptest.NewRecorder()

			b, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/test/%d/%d", tt.taskId, tt.projectId), bytes.NewReader(b))
			require.NoError(t, err)

			router.ServeHTTP(w, req)

			var respBody struct {
				Update bool `json:"updated"`
			}

			require.NoError(t, json.NewDecoder(w.Body).Decode(&respBody))
			require.Equal(t, tt.expReturn, respBody.Update)
			require.Equal(t, tt.expStatusCode, w.Result().StatusCode)
		})
	}
}

func TestRestHandl_Delete(t *testing.T) {
	tests := []struct {
		testName string

		taskId    uint32
		projectId uint32

		expDeleteUC bool
		in          *deletemodel.DeleteTaskInput
		out         *deletemodel.DeleteTaskOutput
		returnErr   error

		expOut        bool
		expStatusCode int
	}{
		{
			testName: "Success",

			taskId:    1,
			projectId: 1,

			expDeleteUC: true,
			in:          deletemodel.NewDeleteTaskInput(1, 1),
			out:         deletemodel.NewDeleteTaskOutput(true),
			returnErr:   nil,

			expOut:        true,
			expStatusCode: http.StatusOK,
		}, {
			testName: "Invalid task id",

			taskId:    0,
			projectId: 1,

			expDeleteUC: false,

			expOut:        false,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Invalid project id",

			taskId:    1,
			projectId: 0,

			expDeleteUC: false,

			expOut:        false,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Task not found",

			taskId:    1,
			projectId: 1,

			expDeleteUC: true,
			in:          deletemodel.NewDeleteTaskInput(1, 1),
			out:         deletemodel.NewDeleteTaskOutput(false),
			returnErr:   deleteerr.ErrTaskNotFound,

			expOut:        false,
			expStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			deleteUC := handlmocks.NewMockDeleteTaskUsecase(ctrl)
			if tt.expDeleteUC {
				deleteUC.EXPECT().Execute(gomock.Any(), tt.in).
					Return(tt.out, tt.returnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			handl := NewRestHandler(log, nil, nil, deleteUC, nil, nil)

			router := gin.New()
			router.DELETE("/test/:task_id/:project_id", handl.Delete)

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/test/%d/%d", tt.taskId, tt.projectId), nil)
			require.NoError(t, err)

			router.ServeHTTP(w, req)

			var respBody struct {
				Deleted bool `json:"deleted"`
			}

			require.NoError(t, json.NewDecoder(w.Body).Decode(&respBody))
			require.Equal(t, tt.expOut, respBody.Deleted)
			require.Equal(t, tt.expStatusCode, w.Result().StatusCode)
		})
	}
}

func TestRestHandl_GetAll(t *testing.T) {
	timeNow := time.Now().Round(1)

	tests := []struct {
		testName string

		projectId uint32

		expGetAllUC bool
		in          *getallmodel.GetAllTasksInput
		out         *getallmodel.GetAllTasksOutput
		returnErr   error

		expOut        []*taskdomain.TaskDomain
		expStatusCode int
	}{
		{
			testName: "Success",

			projectId: 1,

			expGetAllUC: true,
			in:          getallmodel.NewGetAllTasksInput(1),
			out:         getallmodel.NewGetAllTasksOutput([]*taskdomain.TaskDomain{{Id: 1, ProjectId: 1, Description: "asd", Deadline: timeNow}, {Id: 2, ProjectId: 1, Description: "dsa", Deadline: timeNow}}),
			returnErr:   nil,

			expOut:        []*taskdomain.TaskDomain{{Id: 1, ProjectId: 1, Description: "asd", Deadline: timeNow}, {Id: 2, ProjectId: 1, Description: "dsa", Deadline: timeNow}},
			expStatusCode: http.StatusOK,
		}, {
			testName: "Invalid project id",

			projectId: 0,

			expGetAllUC: false,

			expOut:        nil,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Tasks not found",

			projectId: 1,

			expGetAllUC: true,
			in:          getallmodel.NewGetAllTasksInput(1),
			out:         getallmodel.NewGetAllTasksOutput(nil),
			returnErr:   getallerr.ErrTasksNotFound,

			expOut:        nil,
			expStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			getAllUC := handlmocks.NewMockGetAllTasksUsecase(ctrl)
			if tt.expGetAllUC {
				getAllUC.EXPECT().Execute(gomock.Any(), tt.in).
					Return(tt.out, tt.returnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			handl := NewRestHandler(log, nil, nil, nil, getAllUC, nil)

			router := gin.New()
			router.GET("/test/:project_id", handl.GetAll)

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/test/%d", tt.projectId), nil)
			require.NoError(t, err)

			router.ServeHTTP(w, req)

			var respBody struct {
				Tasks []*taskdomain.TaskDomain `json:"tasks"`
			}

			require.NoError(t, json.NewDecoder(w.Body).Decode(&respBody))
			require.Equal(t, tt.expOut, respBody.Tasks)
			require.Equal(t, tt.expStatusCode, w.Result().StatusCode)
		})
	}
}

func TestRestHandl_Get(t *testing.T) {
	timeNow := time.Now().Round(1)

	tests := []struct {
		testName string

		taskId    uint32
		projectId uint32

		expGetAllUC bool
		in          *getmodel.GetTaskInput
		out         *getmodel.GetTaskOutput
		returnErr   error

		expOut        *taskdomain.TaskDomain
		expStatusCode int
	}{
		{
			testName: "Success",

			taskId:    1,
			projectId: 1,

			expGetAllUC: true,
			in:          getmodel.NewGetTaskInput(1, 1),
			out:         getmodel.NewGetTaskOutput(&taskdomain.TaskDomain{Id: 1, ProjectId: 1, Description: "asd", Deadline: timeNow}),
			returnErr:   nil,

			expOut:        &taskdomain.TaskDomain{Id: 1, ProjectId: 1, Description: "asd", Deadline: timeNow},
			expStatusCode: http.StatusOK,
		}, {
			testName: "Invalid project id",

			taskId:    1,
			projectId: 0,

			expGetAllUC: false,

			expOut:        nil,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Tasks not found",

			taskId:    1,
			projectId: 1,

			expGetAllUC: true,
			in:          getmodel.NewGetTaskInput(1, 1),
			out:         getmodel.NewGetTaskOutput(nil),
			returnErr:   geterr.ErrTaskNotFound,

			expOut:        nil,
			expStatusCode: http.StatusNotFound,
		}, {
			testName: "Invalid task id",

			taskId:    0,
			projectId: 1,

			expGetAllUC: false,

			expOut:        nil,
			expStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			getUC := handlmocks.NewMockGetTaskUsecase(ctrl)
			if tt.expGetAllUC {
				getUC.EXPECT().Execute(gomock.Any(), tt.in).
					Return(tt.out, tt.returnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			handl := NewRestHandler(log, nil, nil, nil, nil, getUC)

			router := gin.New()
			router.GET("/test/:project_id/:task_id", handl.Get)

			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/test/%d/%d", tt.projectId, tt.taskId), nil)
			require.NoError(t, err)

			router.ServeHTTP(w, req)

			var respBody struct {
				Tasks *taskdomain.TaskDomain `json:"task"`
			}

			require.NoError(t, json.NewDecoder(w.Body).Decode(&respBody))
			require.Equal(t, tt.expOut, respBody.Tasks)
			require.Equal(t, tt.expStatusCode, w.Result().StatusCode)
		})
	}
}
