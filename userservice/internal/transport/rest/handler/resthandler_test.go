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
	regerr "userservice/internal/usecase/errors/registration"
	regmodel "userservice/internal/usecase/models/registration"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../usecase/interfaces/registration.go -destination=mocks/mock_registration.go -package=handlmocks
func TestNewRestHandler_Registration(t *testing.T) {
	tests := []struct {
		testName   string
		body       []byte
		needExpect bool
		returnData regmodel.RegOutput
		returnErr  error
		expRes     string
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
			returnData: regmodel.RegOutput{IsRegistered: true},
			returnErr:  nil,
			expRes:     "{\"is_registered\":true}",
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
			expRes:     "{\"error\":\"user already exists\"}",
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
			returnData: regmodel.RegOutput{IsRegistered: false},
			returnErr:  nil,
			expRes:     "{\"errors\":{\"FirstName\":\"field is required\"}}",
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
			returnData: regmodel.RegOutput{IsRegistered: false},
			returnErr:  nil,
			expRes:     "{\"errors\":{\"FirstName\":\"field is required\"}}",
			expStatus:  400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			regMock := handlmocks.NewMockRegistrationUsecase(mockCtrl)
			if tt.needExpect {
				regMock.EXPECT().RegUser(gomock.Any(), gomock.Any()).
					Return(&tt.returnData, tt.returnErr)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			handl := NewRestHandler(log, regMock, nil)

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
			assert.Equal(t, tt.expRes, string(data))
		})
	}
}
