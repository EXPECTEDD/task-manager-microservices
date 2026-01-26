package login

import (
	"context"
	"io"
	"log/slog"
	"testing"
	userdomain "userservice/internal/domain/user"
	"userservice/internal/repository/hasher"
	storagerepo "userservice/internal/repository/storage"
	logerr "userservice/internal/usecase/errors/login"
	logmocks "userservice/internal/usecase/implementations/login/mocks"
	logmodel "userservice/internal/usecase/models/login"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../repository/storage/storagerepo.go -destination=mocks/mock_storage.go -package=logmocks
//go:generate mockgen -source=./../../../repository/session/sessionrepo.go -destination=mocks/mock_session.go -package=logmocks
//go:generate mockgen -source=./../../../repository/hasher/password_hasher.go -destination=mocks/mock_password_hasher.go -package=logmocks
//go:generate mockgen -source=./../../../repository/idgenerator/id_generator.go -destination=mocks/mock_id_generator.go -package=logmocks
func TestLogin(t *testing.T) {
	tests := []struct {
		testName string

		expFindByEmail       bool
		findByEmailInput     string
		findByEmailUdReturn  *userdomain.UserDomain
		findByEmailErrReturn error

		expComparePassword           bool
		comparePasswordHashPassInput []byte
		comparePasswordPassInput     []byte
		comparePasswordErrReturn     error

		expSave            bool
		saveSessionIdInput string
		saveUserIdInput    uint32
		saveErrReturn      error

		expNew    bool
		newReturn string

		loginInput     *logmodel.LoginInput
		expLoginOutput *logmodel.LoginOutput
		expLoginErr    error
	}{
		{
			testName: "Success",

			expFindByEmail:   true,
			findByEmailInput: "gmail@gmail.com",
			findByEmailUdReturn: userdomain.NewUserDomain(
				1,
				"Ivan",
				"Ivanovich",
				"Ivanov",
				"hashPass",
				"gmail@gmail.com",
			),
			findByEmailErrReturn: nil,

			expComparePassword:           true,
			comparePasswordHashPassInput: []byte("hashPass"),
			comparePasswordPassInput:     []byte("pass"),
			comparePasswordErrReturn:     nil,

			expSave:            true,
			saveSessionIdInput: "1",
			saveUserIdInput:    1,
			saveErrReturn:      nil,

			expNew:    true,
			newReturn: "1",

			loginInput: logmodel.NewLoginInput("gmail@gmail.com", "pass"),
			expLoginOutput: logmodel.NewLoginOutput(
				"1",
				"Ivan",
				"Ivanovich",
				"Ivanov",
			),
			expLoginErr: nil,
		}, {
			testName: "User not found",

			expFindByEmail:       true,
			findByEmailInput:     "gmail@gmail.com",
			findByEmailUdReturn:  nil,
			findByEmailErrReturn: storagerepo.ErrNoRows,

			expComparePassword: false,
			expSave:            false,
			expNew:             false,

			loginInput:     logmodel.NewLoginInput("gmail@gmail.com", "pass"),
			expLoginOutput: &logmodel.LoginOutput{},
			expLoginErr:    logerr.ErrUserNotFound,
		}, {
			testName: "Wrong password",

			expFindByEmail:   true,
			findByEmailInput: "gmail@gmail.com",
			findByEmailUdReturn: userdomain.NewUserDomain(
				1,
				"Ivan",
				"Ivanovich",
				"Ivanov",
				"hashPass",
				"gmail@gmail.com",
			),
			findByEmailErrReturn: nil,

			expComparePassword:           true,
			comparePasswordHashPassInput: []byte("hashPass"),
			comparePasswordPassInput:     []byte("pass"),
			comparePasswordErrReturn:     hasher.ErrWrongPassword,

			expSave: false,
			expNew:  false,

			loginInput:     logmodel.NewLoginInput("gmail@gmail.com", "pass"),
			expLoginOutput: &logmodel.LoginOutput{},
			expLoginErr:    logerr.ErrWrongPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			storageMock := logmocks.NewMockStorageRepo(ctrl)
			if tt.expFindByEmail {
				storageMock.EXPECT().FindByEmail(gomock.Any(), tt.findByEmailInput).
					Return(tt.findByEmailUdReturn, tt.findByEmailErrReturn)
			}

			passHasherMock := logmocks.NewMockPasswordHasher(ctrl)
			if tt.expComparePassword {
				passHasherMock.EXPECT().ComparePassword(tt.comparePasswordHashPassInput, tt.comparePasswordPassInput).
					Return(tt.comparePasswordErrReturn)
			}

			sessionMock := logmocks.NewMockSessionRepo(ctrl)
			if tt.expSave {
				sessionMock.EXPECT().Save(gomock.Any(), tt.saveSessionIdInput, tt.saveUserIdInput).
					Return(tt.saveErrReturn)
			}

			idgen := logmocks.NewMockIDGenerator(ctrl)
			if tt.expNew {
				idgen.EXPECT().New().Return(tt.newReturn)
			}

			logUC := NewLoginUserUC(log, storageMock, passHasherMock, sessionMock, idgen)
			lo, err := logUC.Execute(context.Background(), tt.loginInput)
			assert.ErrorIs(t, err, tt.expLoginErr)
			assert.Equal(t, tt.expLoginOutput, lo)
		})
	}
}
