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

func TestNotePostgres_Create(t *testing.T) {
	testNote := mother.NoteMother()
	testAccount := mother.AccountMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inNote        model.Note
		inID          int
		expectedError error
	}{
		{
			testName:      "NoteCreatedSuccessfully",
			prepOps:       []string{fmt.Sprintf(testutils.NewAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inNote:        testNote,
			inID:          1,
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
			repo := psql.NewNotePostgres(client, logger)
			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			err = repo.Create(testSuite.inID, &testSuite.inNote)
			assert.Equal(t, testSuite.expectedError, err)

			err = testutils.Cleanup(client, "../../../etc/migrations")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
	err := testutils.CleanupLogs()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotePostgres_GetAll(t *testing.T) {
	testAccount := mother.AccountMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inID          int
		expectedError error
	}{
		{
			testName:      "NotesCollectedSuccessfully",
			prepOps:       []string{fmt.Sprintf(testutils.NewAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:          1,
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
			repo := psql.NewNotePostgres(client, logger)
			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			_, err = repo.GetAll(testSuite.inID)
			assert.Equal(t, testSuite.expectedError, err)

			err = testutils.Cleanup(client, "../../../etc/migrations")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
	err := testutils.CleanupLogs()
	if err != nil {
		t.Fatal(err)
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
		expectedError       error
	}{
		{
			testName:            "NoteCollectedSuccessfully",
			prepOps:             []string{fmt.Sprintf(testutils.NewAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:                1,
			noteShouldBeCreated: true,
			expectedError:       nil,
		},
		{
			testName:            "NotCanNotBeFound",
			prepOps:             []string{fmt.Sprintf(testutils.NewAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:                1,
			noteShouldBeCreated: false,
			expectedError:       e.ClientNoteError,
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

			assert.Equal(t, testSuite.expectedError, err)

			err = testutils.Cleanup(client, "../../../etc/migrations")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
	err := testutils.CleanupLogs()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotePostgres_Delete(t *testing.T) {
	testAccount := mother.AccountMother()
	testNote := mother.NoteMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inID          int
		expectedError error
	}{
		{
			testName:      "NoteDeletedSuccessfully",
			prepOps:       []string{fmt.Sprintf(testutils.NewAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:          1,
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
			repo := psql.NewNotePostgres(client, logger)
			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			err = repo.Create(1, &testNote)

			err = repo.Delete(testSuite.inID, testNote.ID)
			assert.Equal(t, testSuite.expectedError, err)

			err = testutils.Cleanup(client, "../../../etc/migrations")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
	err := testutils.CleanupLogs()
	if err != nil {
		t.Fatal(err)
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
		expectedError       error
	}{
		{
			testName:      "NoteUpdatedSuccessfully",
			prepOps:       []string{fmt.Sprintf(testutils.NewAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:          1,
			inNote:        testNote,
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
			assert.Equal(t, testSuite.expectedError, err)

			err = testutils.Cleanup(client, "../../../etc/migrations")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
	err := testutils.CleanupLogs()
	if err != nil {
		t.Fatal(err)
	}
}
