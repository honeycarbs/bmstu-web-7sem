//go:build unit
// +build unit

package psql_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"neatly/internal/model"
	"neatly/internal/model/mother"
	"neatly/internal/repository/psql"
	"neatly/pkg/logging"
	"neatly/pkg/testutils"
	"testing"
)

func TestTagPostgres_Create(t *testing.T) {
	testAccount := mother.AccountMother()
	testNote := mother.NoteMother()
	testTag := mother.TagMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inID          int
		inNote        model.Note
		inTag         model.Tag
		expectedError error
	}{
		{
			testName:      "TagCreatedSuccessfully",
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
			repo := psql.NewTagPostgres(client, logger)
			noteRepo := psql.NewNotePostgres(client, logger)

			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			err = noteRepo.Create(1, &testNote)
			if err != nil {
				t.Fatalf("Error: %s\n", err)
			}

			err = repo.Create(1, testNote.ID, &testTag)
			logger.Info(err)

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

func TestTagPostgres_Assign(t *testing.T) {
	testAccount := mother.AccountMother()
	testNote := mother.NoteMother()
	testTag := mother.TagMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inID          int
		inNote        model.Note
		inTag         model.Tag
		expectedError error
	}{
		{
			testName:      "TagAssignedSuccessfully",
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
			repo := psql.NewTagPostgres(client, logger)
			noteRepo := psql.NewNotePostgres(client, logger)

			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			err = noteRepo.Create(1, &testNote)
			if err != nil {
				t.Fatalf("Error: %s\n", err)
			}

			err = repo.Create(1, testNote.ID, &testTag)

			err = repo.Assign(testTag.ID, testNote.ID, 1)

			logger.Info(err)

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

func TestTagPostgres_GetAll(t *testing.T) {
	testAccount := mother.AccountMother()
	testNote := mother.NoteMother()
	testTag := mother.TagMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inID          int
		inNote        model.Note
		inTag         model.Tag
		expectedError error
	}{
		{
			testName:      "TagsCollectedSuccessfully",
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
			repo := psql.NewTagPostgres(client, logger)
			noteRepo := psql.NewNotePostgres(client, logger)

			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			err = noteRepo.Create(1, &testNote)
			if err != nil {
				t.Fatalf("Error: %s\n", err)
			}

			err = repo.Create(1, testNote.ID, &testTag)

			_, err = repo.GetAll(1)

			logger.Info(err)

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

func TestTagPostgres_GetAllByNote(t *testing.T) {
	testAccount := mother.AccountMother()
	testNote := mother.NoteMother()
	testTag := mother.TagMother()

	testSuites := []struct {
		testName      string
		prepOps       []string
		inID          int
		inNote        model.Note
		inTag         model.Tag
		expectedError error
	}{
		{
			testName:      "TagsCollectedSuccessfully",
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
			repo := psql.NewTagPostgres(client, logger)
			noteRepo := psql.NewNotePostgres(client, logger)

			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			err = noteRepo.Create(1, &testNote)
			if err != nil {
				t.Fatalf("Error: %s\n", err)
			}

			err = repo.Create(1, testNote.ID, &testTag)

			_, err = repo.GetAllByNote(1, testTag.ID)

			logger.Info(err)

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

func TestTagPostgres_GetOne(t *testing.T) {
	testAccount := mother.AccountMother()
	testNote := mother.NoteMother()
	testTag := mother.TagMother()

	testSuites := []struct {
		testName            string
		prepOps             []string
		inID                int
		inNote              model.Note
		inTag               model.Tag
		noteShouldBeCreated bool
		expectedError       error
	}{
		{
			testName:            "TagCollectedSuccessfully",
			prepOps:             []string{fmt.Sprintf(testutils.NewAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:                1,
			inNote:              testNote,
			noteShouldBeCreated: true,
			expectedError:       nil,
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
			repo := psql.NewTagPostgres(client, logger)
			noteRepo := psql.NewNotePostgres(client, logger)

			for _, op := range testSuite.prepOps {
				_, err = client.DB.Exec(op)
				if err != nil {
					t.Fatalf("sql.Exec: Error: %s\n", err)
				}
			}
			if testSuite.noteShouldBeCreated {
				err = noteRepo.Create(1, &testNote)
				if err != nil {
					t.Fatalf("Error: %s\n", err)
				}
			}

			err = repo.Create(1, testNote.ID, &testTag)
			if testSuite.noteShouldBeCreated {
				err = repo.Assign(testTag.ID, testNote.ID, 1)
			}

			_, err = repo.GetOne(1, testTag.ID)

			logger.Info(err)

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
