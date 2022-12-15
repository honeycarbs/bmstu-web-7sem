//go:build unit
// +build unit

package sqlite_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"neatly/internal/model"
	"neatly/internal/model/mother"
	"neatly/internal/repository/psql"
	"neatly/pkg/e"
	"neatly/pkg/logging"
	"neatly/pkg/testutils"
	"testing"
)

func TestAccountPostgres_CreateAccount(t *testing.T) {
	testAccount := mother.AccountMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inAccount     model.Account
		expectedError error
	}{
		{
			testName:      "AccountCreated",
			prepOps:       []string{},
			inAccount:     testAccount,
			expectedError: nil,
		},
		{
			testName:      "AccountAlreadyExists",
			prepOps:       []string{fmt.Sprintf(testutils.NewAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inAccount:     testAccount,
			expectedError: e.ClientAccountError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := testutils.Setup("../../../etc/migrations")
			if err != nil {
				t.Fatal(err)
			}

			logging.Init()
			logger := logging.GetLogger()
			repo := psql.NewAccountPostgres(client, logger)
			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					logger.Info("Account creating")
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			logger.Info(testSuite.inAccount)
			err = repo.CreateAccount(&testSuite.inAccount)
			logger.Info(err)

			assert.Equal(t, testSuite.expectedError, err)

			err = testutils.Cleanup(client, "../../../etc/migrations")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAccountPostgres_AuthorizeAccount(t *testing.T) {
	testAccount := mother.AccountMother()
	testAccountDoesNotExist := testAccount
	testAccountDoesNotExist.Username = "test2"

	testSuites := []struct {
		testName      string
		prepOps       []string
		inAccount     model.Account
		expectedError error
	}{
		{
			testName:      "AccountDoesNotExist",
			prepOps:       []string{},
			inAccount:     testAccountDoesNotExist,
			expectedError: e.ClientAuthorizeError,
		},
		{
			testName:      "AccountExists",
			prepOps:       []string{fmt.Sprintf(testutils.NewAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inAccount:     testAccount,
			expectedError: nil,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := testutils.Setup("../../../etc/migrations")
			if err != nil {
				t.Fatal(err)
			}

			logging.Init()
			logger := logging.GetLogger()
			repo := psql.NewAccountPostgres(client, logger)
			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			logger.Info(testSuite.inAccount)
			err = repo.AuthorizeAccount(&testSuite.inAccount)

			assert.Equal(t, testSuite.expectedError, err)

			err = testutils.Cleanup(client, "../../../etc/migrations")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
