package psql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"neatly/internal/model"
	"neatly/internal/model/mother"
	"neatly/pkg/client/psqlclient"
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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewAccountPostgres(client, logger)
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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewAccountPostgres(client, logger)
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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewAccountPostgres(client, logger)
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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewNotePostgres(client, logger)
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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewNotePostgres(client, logger)
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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewNotePostgres(client, logger)
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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewNotePostgres(client, logger)
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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewNotePostgres(client, logger)
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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewTagPostgres(client, logger)
			noteRepo := NewNotePostgres(client, logger)

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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewTagPostgres(client, logger)
			noteRepo := NewNotePostgres(client, logger)

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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewTagPostgres(client, logger)
			noteRepo := NewNotePostgres(client, logger)

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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewTagPostgres(client, logger)
			noteRepo := NewNotePostgres(client, logger)

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
			client, err := psqlclient.NewTestClient()
			if err != nil {
				t.Fatal(err)
			}
			defer psqlclient.TestClientClose(client)

			logging.Init()
			logger := logging.GetLogger()
			repo := NewTagPostgres(client, logger)
			noteRepo := NewNotePostgres(client, logger)

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
