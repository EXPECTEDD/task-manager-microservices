package resthandler

import (
	"bytes"
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
	"github.com/stretchr/testify/assert"
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
		expRes     []byte
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
			returnData: regmodel.RegOutput{IsRegistered: true},
			returnErr:  nil,
			expRes:     []byte(`{"is_registered":true}`),
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
			returnData: regmodel.RegOutput{IsRegistered: false},
			returnErr:  regerr.ErrUserAlreadyExists,
			expRes:     []byte(`{"error":"user already exists"}`),
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
			returnData: regmodel.RegOutput{IsRegistered: false},
			returnErr:  nil,
			expRes:     []byte(`{"errors":{"FirstName":"field is required"}}`),
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
			returnData: regmodel.RegOutput{IsRegistered: false},
			returnErr:  nil,
			expRes:     []byte(`{"errors":{"FirstName":"field is required"}}`),
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

			serv := httptest.NewServer(router)
			defer serv.Close()

			rest, err := http.Post(serv.URL+"/test", "application/json", bytes.NewReader(tt.body))
			assert.NoError(t, err)
			defer rest.Body.Close()

			data, err := io.ReadAll(rest.Body)
			assert.Equal(t, tt.expStatus, rest.StatusCode)
			assert.Equal(t, tt.expRes, data)
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

		expBody       []byte
		expStatusCode int
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

			expBody:       []byte(`{"user":{"first_name":"Ivan","middle_name":"Ivanovich","last_name":"Ivanov"}}`),
			expStatusCode: 200,
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

			expBody:       []byte(`{"error":"user not found"}`),
			expStatusCode: 404,
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

			expBody:       []byte(`{"error":"wrong password"}`),
			expStatusCode: 401,
		}, {
			testName:  "Empty field email",
			cookieTTL: time.Duration(3600) * time.Second,

			expectLogin:    false,
			loginOutReturn: &logmodel.LoginOutput{},

			reqBody: []byte(`{
				"password":"somePass"
			}`),

			expBody:       []byte(`{"errors":{"Email":"field is required"}}`),
			expStatusCode: 400,
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

			serv := httptest.NewServer(router)
			defer serv.Close()

			resp, err := http.Post(serv.URL+"/test", "application/json", bytes.NewReader(tt.reqBody))

			assert.NoError(t, err)
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.expStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expBody, body)
		})
	}
}
