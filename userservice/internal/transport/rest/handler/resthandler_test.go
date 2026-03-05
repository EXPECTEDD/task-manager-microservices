package resthandler

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	handlmocks "userservice/internal/transport/rest/handler/mocks"
	"userservice/internal/transport/rest/middleware"
	logerr "userservice/internal/usecase/errors/login"
	regerr "userservice/internal/usecase/errors/registration"
	logmodel "userservice/internal/usecase/models/login"
	regmodel "userservice/internal/usecase/models/registration"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../usecase/interfaces/registration.go -destination=mocks/mock_registration.go -package=handlmocks
func TestRestHandler_Registration(t *testing.T) {
	tests := []struct {
		testName   string
		body       []byte
		needExpect bool
		cookieTTL  time.Duration
		returnData regmodel.RegOutput
		returnErr  error
		expUserId  uint32
		expStatus  int
	}{
		{
			testName: "Success",
			body: []byte(`{
					"first_name":"Ivan",
					"middle_name":"Ivanovich",
					"last_name":"Ivanov",
					"password":"somePass",
					"email":"gmail@gmail.com"
					}`),
			needExpect: true,
			cookieTTL:  time.Duration(3600) * time.Second,
			returnData: regmodel.RegOutput{UserId: 1},
			returnErr:  nil,
			expUserId:  1,
			expStatus:  200,
		}, {
			testName: "User already exists",
			body: []byte(`{
					"first_name":"Ivan",
					"middle_name":"Ivanovich",
					"last_name":"Ivanov",
					"password":"somePass",
					"email":"gmail@gmail.com"
					}`),
			needExpect: true,
			returnData: regmodel.RegOutput{UserId: 0},
			returnErr:  regerr.ErrUserAlreadyExists,
			expUserId:  0,
			expStatus:  409,
		}, {
			testName: "Missing field first_name",
			body: []byte(`{
					"middle_name":"Ivanovich",
					"last_name":"Ivanov",
					"password":"somePass",
					"email":"gmail@gmail.com"
					}`),
			needExpect: false,
			cookieTTL:  time.Duration(3600) * time.Second,
			returnData: regmodel.RegOutput{UserId: 1},
			returnErr:  nil,
			expUserId:  0,
			expStatus:  400,
		}, {
			testName: "Empty field first_name",
			body: []byte(`{
					"first_name":"",
					"middle_name":"Ivanovich",
					"last_name":"Ivanov",
					"password":"somePass",
					"email":"gmail@gmail.com"
					}`),
			needExpect: false,
			cookieTTL:  time.Duration(3600) * time.Second,
			returnData: regmodel.RegOutput{UserId: 0},
			returnErr:  nil,
			expUserId:  0,
			expStatus:  400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			regMock := handlmocks.NewMockRegisterUserUsecase(mockCtrl)
			if tt.needExpect {
				regMock.EXPECT().Execute(gomock.Any(), gomock.Any()).
					Return(&tt.returnData, tt.returnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			handl := NewRestHandler(log, tt.cookieTTL, regMock, nil)

			gin.SetMode(gin.DebugMode)
			router := gin.New()
			router.Use(gin.Recovery())
			router.Use(middleware.TimeoutMiddleware(15 * time.Second))

			router.POST("/test", handl.Registration)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(tt.body))
			require.NoError(t, err)

			router.ServeHTTP(w, req)
			require.NoError(t, err)

			var respBody struct {
				UserId uint32 `json:"user_id"`
			}

			require.NoError(t, json.NewDecoder(w.Body).Decode(&respBody))
			require.Equal(t, tt.expStatus, w.Result().StatusCode)
			require.Equal(t, tt.expUserId, respBody.UserId)
		})
	}
}

//go:generate mockgen -source=./../../../usecase/interfaces/login.go -destination=mocks/mock_login.go -package=handlmocks
func TestRestHandler_Login(t *testing.T) {
	tests := []struct {
		testName  string
		cookieTTL time.Duration

		expectLogin    bool
		loginOutReturn *logmodel.LoginOutput
		loginErrReturn error

		reqBody []byte

		expStatusCode int
		expData       bool
		expFirstName  string
		expMiddleName string
		expLastName   string
	}{
		{
			testName:  "Success",
			cookieTTL: time.Duration(3600) * time.Second,

			expectLogin: true,
			loginOutReturn: logmodel.NewLoginOutput(
				"sessionId",
				"Ivan",
				"Ivanovich",
				"Ivanov",
			),
			loginErrReturn: nil,

			reqBody: []byte(`{
				"email":"gmail@gmail.com",
				"password":"somePass"
			}`),

			expStatusCode: 200,
			expData:       true,
			expFirstName:  "Ivan",
			expMiddleName: "Ivanovich",
			expLastName:   "Ivanov",
		}, {
			testName:  "User not found",
			cookieTTL: time.Duration(3600) * time.Second,

			expectLogin:    true,
			loginOutReturn: &logmodel.LoginOutput{},
			loginErrReturn: logerr.ErrUserNotFound,

			reqBody: []byte(`{
				"email":"gmail@gmail.com",
				"password":"somePass"
			}`),

			expStatusCode: 404,
			expData:       false,
		}, {
			testName:  "Wrong password",
			cookieTTL: time.Duration(3600) * time.Second,

			expectLogin:    true,
			loginOutReturn: &logmodel.LoginOutput{},
			loginErrReturn: logerr.ErrWrongPassword,

			reqBody: []byte(`{
				"email":"gmail@gmail.com",
				"password":"somePass"
			}`),

			expStatusCode: 401,
			expData:       false,
		}, {
			testName:  "Empty field email",
			cookieTTL: time.Duration(3600) * time.Second,

			expectLogin:    false,
			loginOutReturn: &logmodel.LoginOutput{},

			reqBody: []byte(`{
				"password":"somePass"
			}`),

			expStatusCode: 400,
			expData:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loginUCMock := handlmocks.NewMockLoginUserUsecase(ctrl)
			if tt.expectLogin {
				loginUCMock.EXPECT().Execute(gomock.Any(), gomock.Any()).
					Return(tt.loginOutReturn, tt.loginErrReturn)
			}
			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			handl := NewRestHandler(log, tt.cookieTTL, nil, loginUCMock)

			gin.SetMode(gin.DebugMode)
			router := gin.New()
			router.Use(gin.Recovery())
			router.Use(middleware.TimeoutMiddleware(time.Duration(15) * time.Second))

			router.POST("/test", handl.Login)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(tt.reqBody))
			require.NoError(t, err)

			router.ServeHTTP(w, req)
			require.NoError(t, err)

			var respBody struct {
				Data struct {
					FirstName  string `json:"first_name"`
					MiddleName string `json:"middle_name"`
					LastName   string `json:"last_name"`
				} `json:"user"`
			}

			require.NoError(t, json.NewDecoder(w.Body).Decode(&respBody))
			require.Equal(t, tt.expStatusCode, w.Result().StatusCode)
			if tt.expData {
				require.Equal(t, tt.expFirstName, respBody.Data.FirstName)
				require.Equal(t, tt.expMiddleName, respBody.Data.MiddleName)
				require.Equal(t, tt.expLastName, respBody.Data.LastName)
			}
		})
	}
}
