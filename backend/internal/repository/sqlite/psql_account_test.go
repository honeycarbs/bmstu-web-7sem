//go:build unit
// +build unit

package sqlite_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"neatly/internal/model"
	"neatly/internal/model/mother"
	"neatly/internal/repository/psql"
	"neatly/pkg/dbclient"
	"neatly/pkg/logging"
	"testing"
)

var (
	newAccountQuery = "INSERT INTO users (name, username, email, password_hash) VALUES ('%v', '%v', '%v', '%v') RETURNING id"
)

func TestAccountPostgres_CreateAccount(t *testing.T) {
	testAccount := mother.AccountMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inAccount     model.Account
		errorExpected bool
	}{
		{
			testName:      "AccountCreated",
			prepOps:       []string{},
			inAccount:     testAccount,
			errorExpected: false,
		},
		{
			testName:      "AccountAlreadyExists",
			prepOps:       []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inAccount:     testAccount,
			errorExpected: true,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := dbclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer dbclient.TestClientClose(client)

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
			err = repo.CreateAccount(&testSuite.inAccount)
			logger.Info(err)

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
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
		errorExpected bool
	}{
		{
			testName:      "AccountDoesNotExist",
			prepOps:       []string{},
			inAccount:     testAccountDoesNotExist,
			errorExpected: true,
		},
		{
			testName:      "AccountExists",
			prepOps:       []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inAccount:     testAccount,
			errorExpected: false,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := dbclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer dbclient.TestClientClose(client)

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
			logger.Info(err)

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
	}
}

func TestAccountPostgres_GetOne(t *testing.T) {
	testAccount := mother.AccountMother()
	testAccountDoesNotExist := testAccount
	testAccountDoesNotExist.Username = "test2"

	testSuites := []struct {
		testName      string
		prepOps       []string
		inAccount     int
		errorExpected bool
	}{
		{
			testName:      "AccountDoesNotExist",
			prepOps:       []string{},
			inAccount:     0,
			errorExpected: true,
		},
		{
			testName:      "AccountExists",
			prepOps:       []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inAccount:     1,
			errorExpected: false,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := dbclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer dbclient.TestClientClose(client)

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
			_, err = repo.GetOne(testSuite.inAccount)
			logger.Info(err)

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
	}
}
