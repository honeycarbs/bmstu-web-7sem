//go:build unit
// +build unit

package account

import (
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"neatly/internal/model"
	"neatly/internal/model/mother"
	"neatly/internal/repository"
	"neatly/internal/repository/mock"
	"neatly/pkg/e"
	"neatly/pkg/logging"
	"os"
	"testing"
)

func TestService_CreateAccount(t *testing.T) {
	type RepoMockBehaviour func(r *mock.MockAccountRepository, a *model.Account)

	testSuites := []struct {
		testName               string
		inAccount              model.Account
		CreateAccountBehaviour RepoMockBehaviour
		outAccount             model.Account
		ExpectedError          error
	}{
		{
			testName:  "ValidAccountRegistration",
			inAccount: mother.AccountMother(),
			CreateAccountBehaviour: func(r *mock.MockAccountRepository, a *model.Account) {
				r.EXPECT().CreateAccount(a).Return(nil)
			},
			outAccount:    mother.AccountMother(),
			ExpectedError: nil,
		},
		{
			testName:  "UserAlreadyExists",
			inAccount: mother.AccountMother(),
			CreateAccountBehaviour: func(r *mock.MockAccountRepository, a *model.Account) {
				r.EXPECT().CreateAccount(a).Return(e.ClientAccountError)
			},
			outAccount:    mother.AccountMother(),
			ExpectedError: e.ClientAccountError,
		},
	}

	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoMock := mock.NewMockAccountRepository(c)
			testSuite.CreateAccountBehaviour(repoMock, &testSuite.inAccount)

			logging.Init()
			repo := &repository.AccountRepositoryImpl{
				AccountRepository: repoMock,
			}
			mockService := NewService(repo, logging.GetLogger())

			err := mockService.CreateAccount(&testSuite.inAccount)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestService_GenerateJWT(t *testing.T) {
	type RepoMockBehaviour func(r *mock.MockAccountRepository, a *model.Account)
	testAccount := mother.AccountMother()
	testAccountInvalidPassword := testAccount
	testAccountInvalidPassword.Password = "kto prochital tot loh"

	err := os.Setenv("CONF_FILE", "../etc/test.yml")
	if err != nil {
		t.Fatalf("Can't set config path: %s", err)
	}

	testSuites := []struct {
		testName                  string
		inAccount                 model.Account
		AuthorizeAccountBehaviour RepoMockBehaviour
		outAccount                model.Account
		ExpectedError             error
		ExpectedTokenVal          string
	}{
		{
			testName:  "AuthorizeSuccessful",
			inAccount: testAccount,
			AuthorizeAccountBehaviour: func(r *mock.MockAccountRepository, a *model.Account) {
				r.EXPECT().AuthorizeAccount(a).Return(nil)
			},
			outAccount:       testAccount,
			ExpectedError:    nil,
			ExpectedTokenVal: mother.TokenMother(),
		},
		{
			testName:  "PasswordDoesNotMatch",
			inAccount: testAccountInvalidPassword,
			AuthorizeAccountBehaviour: func(r *mock.MockAccountRepository, a *model.Account) {
				r.EXPECT().AuthorizeAccount(a).Return(nil)
			},
			outAccount:       testAccount,
			ExpectedError:    errors.New("password does not match"),
			ExpectedTokenVal: "",
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoMock := mock.NewMockAccountRepository(c)
			testSuite.AuthorizeAccountBehaviour(repoMock, &testSuite.inAccount)

			logging.Init()
			logger := logging.GetLogger()
			repo := &repository.AccountRepositoryImpl{
				AccountRepository: repoMock,
			}
			mockService := NewService(repo, logger)

			token, err := mockService.GenerateJWT(&testSuite.inAccount)

			assert.Equal(t, testSuite.ExpectedError, err)
			assert.Equal(t, testSuite.ExpectedTokenVal, token)
		})
	}
}
