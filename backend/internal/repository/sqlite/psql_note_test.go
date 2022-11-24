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

func TestNotePostgres_Create(t *testing.T) {
	testNote := mother.NoteMother()
	testAccount := mother.AccountMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inNote        model.Note
		inID          int
		errorExpected bool
	}{
		{
			testName:      "NoteCreatedSuccessfully",
			prepOps:       []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inNote:        testNote,
			inID:          1,
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
			repo := psql.NewNotePostgres(client, logger)
			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			err = repo.Create(testSuite.inID, &testSuite.inNote)
			logger.Info(err)

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
	}
}

func TestNotePostgres_GetAll(t *testing.T) {
	testAccount := mother.AccountMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inID          int
		errorExpected bool
	}{
		{
			testName:      "NotesCollectedSuccessfully",
			prepOps:       []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:          1,
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
			repo := psql.NewNotePostgres(client, logger)
			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			_, err = repo.GetAll(testSuite.inID)
			logger.Info(err)

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
	}
}

func TestNotePostgres_GetOne(t *testing.T) {
	testAccount := mother.AccountMother()
	testNote := mother.NoteMother()

	testSuites := []struct {
		testName            string
		prepOps             []string
		inID                int
		noteShouldBeCreated bool
		errorExpected       bool
	}{
		{
			testName:            "NoteCollectedSuccessfully",
			prepOps:             []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:                1,
			noteShouldBeCreated: true,
			errorExpected:       false,
		},
		{
			testName:            "NotCanNotBeFound",
			prepOps:             []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:                1,
			noteShouldBeCreated: false,
			errorExpected:       true,
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
			repo := psql.NewNotePostgres(client, logger)
			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			if testSuite.noteShouldBeCreated {
				err = repo.Create(1, &testNote)
			}
			_, err = repo.GetOne(testSuite.inID, testNote.ID)
			logger.Info(err)

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
	}
}

func TestNotePostgres_Delete(t *testing.T) {
	testAccount := mother.AccountMother()
	testNote := mother.NoteMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inID          int
		errorExpected bool
	}{
		{
			testName:      "NoteDeletedSuccessfully",
			prepOps:       []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:          1,
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
			repo := psql.NewNotePostgres(client, logger)
			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			err = repo.Create(1, &testNote)

			err = repo.Delete(testSuite.inID, testNote.ID)
			logger.Info(err)

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
	}
}

func TestNotePostgres_Update(t *testing.T) {
	testAccount := mother.AccountMother()
	testNote := mother.NoteMother()

	testSuites := []struct {
		testName            string
		prepOps             []string
		inID                int
		inNote              model.Note
		noteShouldBeCreated bool
		errorExpected       bool
	}{
		{
			testName:            "NoteUpdatedSuccessfully",
			prepOps:             []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:                1,
			inNote:              testNote,
			noteShouldBeCreated: true,
			errorExpected:       false,
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
			repo := psql.NewNotePostgres(client, logger)
			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			if testSuite.noteShouldBeCreated {
				err = repo.Create(1, &testNote)
			}
			err = repo.Update(testSuite.inID, testSuite.inNote)
			logger.Info(err)

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
	}
}
