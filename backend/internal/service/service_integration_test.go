package service

import (
	"github.com/go-playground/assert/v2"
	"neatly/internal/model"
	"neatly/internal/model/mother"
	"neatly/internal/repository"
	"neatly/internal/service/account"
	"neatly/pkg/e"
	"neatly/pkg/integration"
	"neatly/pkg/logging"
	"testing"
)

func prepareAccountRepo(action string, repo repository.AccountRepository) error {
	switch action {
	case "INSERT":
		{
			a := mother.AccountMother()
			return repo.CreateAccount(&a)
		}
	default:
		return nil
	}
}

func TestAccountCreate(t *testing.T) {
	testSuites := []struct {
		testName      string
		inAccount     model.Account
		outAccount    model.Account
		prepareAction string
		ExpectedError error
	}{
		{
			testName:      "ValidAccountRegistration",
			inAccount:     mother.AccountMother(),
			outAccount:    mother.AccountMother(),
			prepareAction: "DO NOTHING",
			ExpectedError: nil,
		},
		{
			testName:      "UserAlreadyExists",
			inAccount:     mother.AccountMother(),
			outAccount:    mother.AccountMother(),
			prepareAction: "INSERT",
			ExpectedError: e.ClientAccountError,
		},
	}

	logging.Init()
	logger := logging.GetLogger()

	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {

			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			repo := repository.NewAccountRepositoryImpl(client, logger)
			err = prepareAccountRepo(testSuite.prepareAction, repo)
			if err != nil {
				t.Fatalf("Can't do pre-test action: %s", err)
			}

			service := account.NewService(repo, logger)

			err = service.CreateAccount(&testSuite.inAccount)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}
