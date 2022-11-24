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
		errorExpected bool
	}{
		{
			testName:      "TagCreatedSuccessfully",
			prepOps:       []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:          1,
			inNote:        testNote,
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
			//err = repo.Update(testSuite.inID, testSuite.inNote)
			logger.Info(err)

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
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
		errorExpected bool
	}{
		{
			testName:      "TagAssignedSuccessfully",
			prepOps:       []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:          1,
			inNote:        testNote,
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

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
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
		errorExpected bool
	}{
		{
			testName:      "TagsCollectedSuccessfully",
			prepOps:       []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:          1,
			inNote:        testNote,
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

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
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
		errorExpected bool
	}{
		{
			testName:      "TagsCollectedSuccessfully",
			prepOps:       []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:          1,
			inNote:        testNote,
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

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
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
		errorExpected       bool
	}{
		{
			testName:            "TagCollectedSuccessfully",
			prepOps:             []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:                1,
			inNote:              testNote,
			noteShouldBeCreated: true,
			errorExpected:       false,
		},
		{
			testName:            "TagNotFound",
			prepOps:             []string{fmt.Sprintf(newAccountQuery, testAccount.Name, testAccount.Username, testAccount.Email, testAccount.PasswordHash)},
			inID:                1,
			inNote:              testNote,
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

			if testSuite.errorExpected {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, nil, err)
			}
		})
	}
}
