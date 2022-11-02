package account

import (
	"database/sql"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"log"
	"neatly/internal/model"
	"neatly/internal/repository"
	"neatly/internal/repository/mock"
	"neatly/pkg/e"
	"neatly/pkg/jwt"
	"neatly/pkg/logging"
	"os"
	"testing"
)

func AccountMother() model.Account {

	testHash, _ := model.GeneratePasswordHash("testtest")

	return model.Account{
		ID:           0,
		Name:         "Test",
		Username:     "TestTest",
		Email:        "test",
		Password:     "testtest",
		PasswordHash: testHash,
	}
}

func TokenMother() string {
	a := AccountMother()

	logging.Init()
	os.Setenv("CONF_FILE", "../../../etc/config/local.yml")

	token, err := jwt.GenerateAccessToken(a.ID)
	if err != nil {
		log.Fatal("can't create test token")
	}

	return token
}

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
			inAccount: AccountMother(),
			CreateAccountBehaviour: func(r *mock.MockAccountRepository, a *model.Account) {
				r.EXPECT().CreateAccount(a).Return(nil)
			},
			outAccount:    AccountMother(),
			ExpectedError: nil,
		},
		{
			testName:  "UserAlreadyExists",
			inAccount: AccountMother(),
			CreateAccountBehaviour: func(r *mock.MockAccountRepository, a *model.Account) {
				r.EXPECT().CreateAccount(a).Return(e.ClientAccountError)
			},
			outAccount:    AccountMother(),
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
	testAccount := AccountMother()
	testAccountInvalidPassword := testAccount
	testAccountInvalidPassword.Password = "kto prochital tot loh"

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
			ExpectedTokenVal: TokenMother(),
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
func TestService_GetOne(t *testing.T) {
	type RepoMockBehaviour func(r *mock.MockAccountRepository, id int)
	testAccount := AccountMother()

	testSuites := []struct {
		testName            string
		inID                int
		GetAccountBehaviour RepoMockBehaviour
		outAccount          model.Account
		ExpectedError       error
	}{
		{
			testName: "AccountFound",
			inID:     0,
			GetAccountBehaviour: func(r *mock.MockAccountRepository, id int) {
				r.EXPECT().GetOne(id).Return(testAccount, nil)
			},
			outAccount:    testAccount,
			ExpectedError: nil,
		},
		{
			testName: "AccountNotFound",
			inID:     0,
			GetAccountBehaviour: func(r *mock.MockAccountRepository, id int) {
				r.EXPECT().GetOne(id).Return(testAccount, sql.ErrNoRows)
			},
			outAccount:    testAccount,
			ExpectedError: sql.ErrNoRows,
		},
	}

	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoMock := mock.NewMockAccountRepository(c)
			testSuite.GetAccountBehaviour(repoMock, testSuite.inID)

			logging.Init()
			repo := &repository.AccountRepositoryImpl{
				AccountRepository: repoMock,
			}
			mockService := NewService(repo, logging.GetLogger())
			got, err := mockService.GetOne(testSuite.inID)

			assert.Equal(t, testSuite.ExpectedError, err)
			assert.Equal(t, testSuite.outAccount, got)
		})
	}
}
