package registration

import (
	"context"
	"io"
	"log/slog"
	"testing"
	userdomain "userservice/internal/domain/user"
	storagerepo "userservice/internal/repository/storage"
	regerr "userservice/internal/usecase/errors/registration"
	regmocks "userservice/internal/usecase/implementations/registration/mocks"
	regmodel "userservice/internal/usecase/models/registration"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=./../../../repository/storage/storagerepo.go -destination=mocks/mock_storage.go -package=regmocks
//go:generate mockgen -source=./../../../repository/hasher/password_hasher.go -destination=mocks/mock_hasher.go -package=regmocks
func TestRegUser(t *testing.T) {
	tests := []struct {
		testName string

		expectFindByEmail     bool
		findByEmailInput      string
		findByEmailUserReturn *userdomain.UserDomain
		findByEmailErrReturn  error

		expectHash     bool
		hashInput      []byte
		hashPassReturn []byte
		hashErrReturn  error

		expectSave    bool
		saveInput     *userdomain.UserDomain
		saveIdReturn  uint32
		saveErrReturn error

		regUserInput        regmodel.RegInput
		regUserExpectOutput regmodel.RegOutput
		regUserExpectErr    error
	}{
		{
			testName: "Success",

			expectFindByEmail:     true,
			findByEmailInput:      "gmail@gmail.com",
			findByEmailUserReturn: nil,
			findByEmailErrReturn:  storagerepo.ErrNoRows,

			expectHash:     true,
			hashInput:      []byte("somePass"),
			hashPassReturn: []byte("hashPass"),
			hashErrReturn:  nil,

			expectSave: true,
			saveInput: userdomain.NewUserDomain(
				0,
				"Ivan",
				"Ivanovich",
				"Ivanov",
				"hashPass",
				"gmail@gmail.com",
			),
			saveIdReturn:  1,
			saveErrReturn: nil,

			regUserInput: *regmodel.NewRegInput(
				"Ivan",
				"Ivanovich",
				"Ivanov",
				"somePass",
				"gmail@gmail.com",
			),
			regUserExpectOutput: *regmodel.NewRegOutput(
				true,
			),
			regUserExpectErr: nil,
		}, {
			testName: "User already exists",

			expectFindByEmail: true,
			findByEmailInput:  "gmail@gmail.com",
			findByEmailUserReturn: userdomain.NewUserDomain(
				1,
				"Ivan",
				"Ivanovich",
				"Ivanov",
				"somePass",
				"gmail@gmail.com",
			),
			findByEmailErrReturn: nil,

			expectHash: false,
			expectSave: false,

			regUserInput: *regmodel.NewRegInput(
				"Ivan",
				"Ivanovich",
				"Ivanov",
				"somePass",
				"gmail@gmail.com",
			),
			regUserExpectOutput: *regmodel.NewRegOutput(
				false,
			),
			regUserExpectErr: regerr.ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storMock := regmocks.NewMockStorageRepo(ctrl)
			if tt.expectFindByEmail {
				storMock.EXPECT().FindByEmail(gomock.Any(), tt.findByEmailInput).
					Return(tt.findByEmailUserReturn, tt.findByEmailErrReturn)
			}
			if tt.expectSave {
				storMock.EXPECT().Save(gomock.Any(), tt.saveInput).
					Return(tt.saveIdReturn, tt.saveErrReturn)
			}

			hasherMock := regmocks.NewMockPasswordHasher(ctrl)
			if tt.expectHash {
				hasherMock.EXPECT().Hash(tt.hashInput).
					Return(tt.hashPassReturn, tt.hashErrReturn)
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			regUC := NewRegUserUC(log, storMock, hasherMock)

			out, err := regUC.Execute(context.Background(), &tt.regUserInput)
			assert.ErrorIs(t, tt.regUserExpectErr, err)
			assert.Equal(t, tt.regUserExpectOutput.IsRegistered, out.IsRegistered)
		})
	}
}
