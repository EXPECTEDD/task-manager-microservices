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
	deleteerr "projectservice/internal/usecase/error/deleteproject"
	getallerr "projectservice/internal/usecase/error/getallprojects"
	createmodel "projectservice/internal/usecase/models/createproject"
	deletemodel "projectservice/internal/usecase/models/deleteproject"
	getallmodel "projectservice/internal/usecase/models/getallprojects"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../repository/sessionvalidator/session_validator.go -destination=./mocks/mock_session_validator.go -package=resthandlmocks
//go:generate mockgen -source=./../../../usecase/interfaces/create_project.go -destination=./mocks/mock_create_project.go -package=resthandlmocks
//go:generate mockgen -source=./../../../usecase/interfaces/delete_project.go -destination=./mocks/mock_delete_project.go -package=resthandlmocks
//go:generate mockgen -source=./../../../usecase/interfaces/get_all_projects.go -destination=./mocks/mock_get_all_projects.go -package=resthandlmocks
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

		expResp       uint32
		expStatusCode int
	}{
		{
			testName: "Success",

			sessionId: "sessionId",
			userId:    1,

			expCreate:       true,
			createInput:     createmodel.NewCreateProjectInput(1, "Name"),
			createOutput:    createmodel.NewCreateProjectOutput(1),
			createReturnErr: nil,

			contentType: "application/json",
			body: map[string]string{
				"name": "Name",
			},

			expResp:       1,
			expStatusCode: http.StatusOK,
		}, {
			testName: "Missing field name",

			sessionId: "sessionId",
			userId:    1,

			expCreate:       false,
			createInput:     createmodel.NewCreateProjectInput(1, "Name"),
			createOutput:    createmodel.NewCreateProjectOutput(1),
			createReturnErr: nil,

			contentType: "application/json",
			body: map[string]string{
				"nme": "Name",
			},

			expResp:       0,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Invalid name",

			sessionId: "sessionId",
			userId:    1,

			expCreate:       true,
			createInput:     createmodel.NewCreateProjectInput(1, strings.Repeat("Name", 300)),
			createOutput:    createmodel.NewCreateProjectOutput(0),
			createReturnErr: projectdomain.ErrInvalidName,

			contentType: "application/json",
			body: map[string]string{
				"name": strings.Repeat("Name", 300),
			},

			expResp:       0,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Already exists",

			sessionId: "sessionId",
			userId:    1,

			expCreate:       true,
			createInput:     createmodel.NewCreateProjectInput(1, "Name"),
			createOutput:    createmodel.NewCreateProjectOutput(0),
			createReturnErr: createerr.ErrAlreadyExists,

			contentType: "application/json",
			body: map[string]string{
				"name": "Name",
			},

			expResp:       0,
			expStatusCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			createUCMock := resthandlmocks.NewMockCreateProjectUsecase(ctrl)
			if tt.expCreate {
				createUCMock.EXPECT().Execute(gomock.Any(), tt.createInput).
					Return(tt.createOutput, tt.createReturnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			client := resthandlmocks.NewMockSessionValidator(ctrl)
			client.EXPECT().GetIdBySession(gomock.Any(), tt.sessionId).
				Return(tt.userId, nil)

			handl := NewHandler(log, createUCMock, nil, nil)

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
				ProjectId uint32 `json:"project_id"`
			}

			assert.NoError(t, json.NewDecoder(w.Body).Decode(&respBody))
			assert.Equal(t, tt.expResp, respBody.ProjectId)
			assert.Equal(t, tt.expStatusCode, w.Result().StatusCode)
		})
	}
}

func TestRestHandler_Delete(t *testing.T) {
	tests := []struct {
		testName string

		sessionId string
		userId    uint32

		expDelete         bool
		deleteUCInput     *deletemodel.DeleteProjectInput
		deleteUCOutput    *deletemodel.DeleteProjectOutput
		deleteUCReturnErr error

		clientReturnErr error

		body map[string]uint32

		expRespBody   bool
		expStatusCode int
	}{
		{
			testName: "Success",

			sessionId: "sessionId",
			userId:    1,

			expDelete:         true,
			deleteUCInput:     deletemodel.NewDeleteProjectInput(1),
			deleteUCOutput:    deletemodel.NewDeleteProjectOutput(true),
			deleteUCReturnErr: nil,

			clientReturnErr: nil,

			body: map[string]uint32{
				"project_id": 1,
			},

			expRespBody:   true,
			expStatusCode: http.StatusOK,
		}, {
			testName: "Invalid project id",

			sessionId: "sessionId",
			userId:    1,

			expDelete:         false,
			deleteUCInput:     deletemodel.NewDeleteProjectInput(0),
			deleteUCOutput:    deletemodel.NewDeleteProjectOutput(false),
			deleteUCReturnErr: deleteerr.ErrInvalidProjectId,

			clientReturnErr: nil,

			body: map[string]uint32{
				"project_id": 0,
			},

			expRespBody:   false,
			expStatusCode: http.StatusBadRequest,
		}, {
			testName: "Not found",

			sessionId: "sessionId",
			userId:    1,

			expDelete:         true,
			deleteUCInput:     deletemodel.NewDeleteProjectInput(1),
			deleteUCOutput:    deletemodel.NewDeleteProjectOutput(true),
			deleteUCReturnErr: deleteerr.ErrProjectNotFound,

			clientReturnErr: nil,

			body: map[string]uint32{
				"project_id": 1,
			},

			expRespBody:   false,
			expStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			deleteUCMock := resthandlmocks.NewMockDeleteProjectUsecase(ctrl)
			if tt.expDelete {
				deleteUCMock.EXPECT().Execute(gomock.Any(), tt.deleteUCInput).
					Return(tt.deleteUCOutput, tt.deleteUCReturnErr)
			}

			handl := NewHandler(log, nil, deleteUCMock, nil)

			client := resthandlmocks.NewMockSessionValidator(ctrl)

			client.EXPECT().GetIdBySession(gomock.Any(), tt.sessionId).
				Return(tt.userId, tt.clientReturnErr)

			router := gin.New()
			router.Use(middleware.GetSessionMiddleware(log))
			router.Use(middleware.SessionAuthMiddleware(log, client, 10*time.Second))

			router.DELETE("/test", handl.Delete)

			b, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodDelete, "/test", bytes.NewReader(b))

			c := &http.Cookie{
				Name:  "sessionId",
				Value: tt.sessionId,
			}
			req.AddCookie(c)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			var respBody struct {
				IsDeleted bool `json:"is_deleted"`
			}

			assert.NoError(t, json.NewDecoder(w.Body).Decode(&respBody))
			assert.Equal(t, tt.expRespBody, respBody.IsDeleted)
			assert.Equal(t, tt.expStatusCode, w.Result().StatusCode)
		})
	}
}

func TestRestHandler_GetAll(t *testing.T) {
	timeNow := time.Now().Round(0)

	tests := []struct {
		testName string

		userId    uint32
		sessionId string

		ucInput     *getallmodel.GetAllProjectsInput
		ucOutput    *getallmodel.GetAllProjectsOutput
		ucReturnErr error

		expBody       []*projectdomain.ProjectDomain
		expStatusCode int
	}{
		{
			testName: "Success",

			userId:    1,
			sessionId: "sessionId",

			ucInput: getallmodel.NewGetAllProjectsInput(1),
			ucOutput: getallmodel.NewGetAllProjectsOutput([]*projectdomain.ProjectDomain{
				&projectdomain.ProjectDomain{Id: 1, OwnerId: 1, Name: "A", CreatedAt: timeNow},
			}),
			ucReturnErr: nil,

			expBody: []*projectdomain.ProjectDomain{
				&projectdomain.ProjectDomain{Id: 1, OwnerId: 1, Name: "A", CreatedAt: timeNow},
			},
			expStatusCode: http.StatusOK,
		}, {
			testName: "Project not found",

			userId:    1,
			sessionId: "sessionId",

			ucInput:     getallmodel.NewGetAllProjectsInput(1),
			ucOutput:    getallmodel.NewGetAllProjectsOutput([]*projectdomain.ProjectDomain{nil}),
			ucReturnErr: getallerr.ErrProjectsNotFound,

			expBody:       nil,
			expStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			getAllMock := resthandlmocks.NewMockGetAllProjectsUsecase(ctrl)
			getAllMock.EXPECT().Execute(gomock.Any(), tt.ucInput).
				Return(tt.ucOutput, tt.ucReturnErr)

			handl := NewHandler(log, nil, nil, getAllMock)

			client := resthandlmocks.NewMockSessionValidator(ctrl)
			client.EXPECT().GetIdBySession(gomock.Any(), tt.sessionId).
				Return(tt.userId, nil)

			router := gin.New()
			router.Use(middleware.GetSessionMiddleware(log))
			router.Use(middleware.SessionAuthMiddleware(log, client, 10*time.Second))
			router.GET("/test", handl.GetAll)

			req, err := http.NewRequest(http.MethodGet, "/test", nil)
			assert.NoError(t, err)

			c := &http.Cookie{
				Name:  "sessionId",
				Value: tt.sessionId,
			}

			req.AddCookie(c)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			var respBody struct {
				Projects []*projectdomain.ProjectDomain `json:"projects"`
			}

			assert.NoError(t, json.NewDecoder(w.Body).Decode(&respBody))
			assert.Equal(t, tt.expBody, respBody.Projects)
			assert.Equal(t, tt.expStatusCode, w.Result().StatusCode)
		})
	}
}
