package service_test

import (
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"neatly/internal/model"
	"neatly/internal/model/mother"
	"neatly/internal/repository"
	"neatly/internal/service/account"
	"neatly/internal/service/note"
	"neatly/internal/service/tag"
	"neatly/pkg/e"
	"neatly/pkg/integration"
	"neatly/pkg/logging"
	"os"
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

func prepareNoteRepo(action string, userID int, repo repository.NoteRepository) error {
	switch action {
	case "INSERT":
		{
			n := mother.NoteMother()
			return repo.Create(userID, &n)
		}
	case "INSERT x2":
		{
			n := mother.NoteMother()
			err := repo.Create(userID, &n)
			if err != nil {
				return err
			}
			return repo.Create(userID, &n)
		}
	default:
		return nil
	}
}

func prepareTagRepo(action string, userID, noteID int, repo repository.TagRepository) error {
	switch action {
	case "CREATE":
		{
			t := mother.TagMother()
			return repo.Create(userID, noteID, &t)
		}
	case "ASSIGN":
		{
			t := mother.TagMother()
			err := repo.Create(userID, noteID, &t)
			if err != nil {
				return err
			}
			return repo.Assign(t.ID, noteID, userID)
		}
	case "ASSIGN WITH LABEL":
		{
			t := mother.TagMother()
			t.Label = "test"
			err := repo.Create(userID, noteID, &t)
			if err != nil {
				return err
			}
			return repo.Assign(t.ID, noteID, userID)
		}
	default:
		return nil
	}
}

func TestIntegration_AccountCreate(t *testing.T) {
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

func TestIntegration_GenerateJWT(t *testing.T) {
	logging.Init()
	logger := logging.GetLogger()

	testAccount := mother.AccountMother()

	testAccountInvalidPassword := testAccount
	testAccountInvalidPassword.Password = "kto prochital tot loh"

	testSuites := []struct {
		testName      string
		inAccount     model.Account
		outAccount    model.Account
		prepareAction string
		ExpectedError error
	}{
		{
			testName:      "AuthorizeSuccessful",
			inAccount:     testAccount,
			outAccount:    testAccount,
			prepareAction: "INSERT",
			ExpectedError: nil,
		},
		{
			testName:      "PasswordDoesNotMatch",
			inAccount:     testAccountInvalidPassword,
			outAccount:    testAccount,
			prepareAction: "INSERT",
			ExpectedError: errors.New("password does not match"),
		},
	}

	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			err := os.Setenv("CONF_FILE", "etc/test.yml")
			if err != nil {
				t.Fatalf("Can't set config path: %s", err)
			}

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

			token, err := service.GenerateJWT(&testSuite.inAccount)
			logger.Info(token)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestIntegration_AccountGetOne(t *testing.T) {
	testAccount := mother.AccountMother()

	testSuites := []struct {
		testName      string
		inID          int
		prepareAction string
		outAccount    model.Account
		ExpectedError error
	}{
		{
			testName:      "AccountFound",
			inID:          1,
			outAccount:    testAccount,
			prepareAction: "INSERT",
			ExpectedError: nil,
		},
		{
			testName:      "AccountNotFound",
			inID:          0,
			outAccount:    testAccount,
			ExpectedError: sql.ErrNoRows,
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

			_, err = service.GetOne(testSuite.inID)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestIntegration_NoteCreate(t *testing.T) {
	testNote := mother.NoteMother()

	logging.Init()
	logger := logging.GetLogger()

	testSuites := []struct {
		testName      string
		inNote        model.Note
		inID          int
		ExpectedError error
	}{
		{
			testName:      "NoteCreatedSuccessfully",
			inNote:        testNote,
			inID:          1,
			ExpectedError: nil,
		},
		{
			testName:      "UserDoesNotExist",
			inNote:        testNote,
			inID:          0,
			ExpectedError: e.InternalDBError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			nr := repository.NewNoteRepositoryImpl(client, logger)
			ar := repository.NewAccountRepositoryImpl(client, logger)
			tr := repository.NewTagRepositoryImpl(client, logger)

			err = prepareAccountRepo("INSERT", ar)
			if err != nil {
				t.Fatalf("Can't do pre-test action: %s", err)
			}

			service := note.NewService(nr, tr, logger)

			err = service.Create(testSuite.inID, &testNote)
			//
			//if testSuite

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestNoteGetAll(t *testing.T) {
	testNote := mother.NoteMother()

	logging.Init()
	logger := logging.GetLogger()

	testSuites := []struct {
		testName             string
		inNote               model.Note
		prepareAccountAction string
		prepareNoteAction    string
		prepareTagAction     string
		inID                 int
		ExpectedError        error
	}{
		{
			testName:             "UserHasNoNotes",
			inNote:               testNote,
			inID:                 1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "DO NOTHING",
			prepareTagAction:     "DO NOTHING",
			ExpectedError:        nil,
		},
		{
			testName:             "UserHasNoteWithoutTags",
			inNote:               testNote,
			inID:                 1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "DO NOTHING",
			ExpectedError:        nil,
		},
		{
			testName:             "UserHasNoteWithTags",
			inNote:               testNote,
			inID:                 1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "ASSIGN",
			ExpectedError:        nil,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			nr := repository.NewNoteRepositoryImpl(client, logger)
			ar := repository.NewAccountRepositoryImpl(client, logger)
			tr := repository.NewTagRepositoryImpl(client, logger)

			err = prepareAccountRepo(testSuite.prepareAccountAction, ar)
			if err != nil {
				t.Fatalf("Can't do account pre-test action: %s", err)
			}

			err = prepareNoteRepo(testSuite.prepareNoteAction, 1, nr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			err = prepareTagRepo(testSuite.prepareTagAction, 1, 1, tr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			service := note.NewService(nr, tr, logger)

			_, err = service.GetAll(testSuite.inID)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestIntegration_NoteGetOne(t *testing.T) {
	testNote := mother.NoteMother()

	logging.Init()
	logger := logging.GetLogger()

	testSuites := []struct {
		testName             string
		inNote               model.Note
		prepareAccountAction string
		prepareNoteAction    string
		prepareTagAction     string
		inID                 int
		ExpectedError        error
	}{
		{
			testName:             "NoteFound",
			inNote:               testNote,
			inID:                 1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "DO NOTHING",
			ExpectedError:        nil,
		},
		{
			testName:             "NoteDoesNotExists",
			inNote:               testNote,
			inID:                 1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "DO NOTHING",
			prepareTagAction:     "DO NOTHING",
			ExpectedError:        e.ClientNoteError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			nr := repository.NewNoteRepositoryImpl(client, logger)
			ar := repository.NewAccountRepositoryImpl(client, logger)
			tr := repository.NewTagRepositoryImpl(client, logger)

			err = prepareAccountRepo(testSuite.prepareAccountAction, ar)
			if err != nil {
				t.Fatalf("Can't do account pre-test action: %s", err)
			}

			err = prepareNoteRepo(testSuite.prepareNoteAction, 1, nr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			err = prepareTagRepo(testSuite.prepareTagAction, 1, 1, tr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			service := note.NewService(nr, tr, logger)

			_, err = service.GetOne(testSuite.inID, 1)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestIntegration_NoteFindByTag(t *testing.T) {
	testNote := mother.NoteMother()

	logging.Init()
	logger := logging.GetLogger()

	testSuites := []struct {
		testName             string
		inNote               model.Note
		prepareAccountAction string
		prepareNoteAction    string
		prepareTagAction     string
		inID                 int
		ExpectedError        error
	}{
		{
			testName:             "NoteFound",
			inNote:               testNote,
			inID:                 1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "DO NOTHING",
			ExpectedError:        nil,
		},
		{
			testName:             "NoteDoesNotExists",
			inNote:               testNote,
			inID:                 1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "DO NOTHING",
			prepareTagAction:     "DO NOTHING",
			ExpectedError:        nil,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			nr := repository.NewNoteRepositoryImpl(client, logger)
			ar := repository.NewAccountRepositoryImpl(client, logger)
			tr := repository.NewTagRepositoryImpl(client, logger)

			err = prepareAccountRepo(testSuite.prepareAccountAction, ar)
			if err != nil {
				t.Fatalf("Can't do account pre-test action: %s", err)
			}

			err = prepareNoteRepo(testSuite.prepareNoteAction, 1, nr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			err = prepareTagRepo(testSuite.prepareTagAction, 1, 1, tr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			service := note.NewService(nr, tr, logger)

			_, err = service.FindByTags(1, []string{"test"})

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestIntegration_NoteUpdate(t *testing.T) {
	testNote := mother.NoteMother()

	logging.Init()
	logger := logging.GetLogger()

	testSuites := []struct {
		testName             string
		inNote               model.Note
		prepareAccountAction string
		prepareNoteAction    string
		prepareTagAction     string
		inID                 int
		ExpectedError        error
	}{
		{
			testName:             "NoteFound",
			inNote:               testNote,
			inID:                 1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "DO NOTHING",
			ExpectedError:        nil,
		},
		{
			testName:             "NoteDoesNotExists",
			inNote:               testNote,
			inID:                 0,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "DO NOTHING",
			prepareTagAction:     "DO NOTHING",
			ExpectedError:        e.ClientNoteError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			nr := repository.NewNoteRepositoryImpl(client, logger)
			ar := repository.NewAccountRepositoryImpl(client, logger)
			tr := repository.NewTagRepositoryImpl(client, logger)

			err = prepareAccountRepo(testSuite.prepareAccountAction, ar)
			if err != nil {
				t.Fatalf("Can't do account pre-test action: %s", err)
			}

			err = prepareNoteRepo(testSuite.prepareNoteAction, 1, nr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			err = prepareTagRepo(testSuite.prepareTagAction, 1, 1, tr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			service := note.NewService(nr, tr, logger)

			un := mother.NoteMother()
			un.ID = testSuite.inID

			err = service.Update(1, un, false)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestIntegration_TagCreate(t *testing.T) {

	testTag := mother.TagMother()

	logging.Init()
	logger := logging.GetLogger()

	testSuites := []struct {
		testName             string
		prepareAccountAction string
		prepareNoteAction    string
		prepareTagAction     string
		inUserID             int
		inNoteID             int
		inTag                model.Tag
		ExpectedError        error
	}{
		{
			testName:             "TagAssigned",
			inUserID:             1,
			inNoteID:             1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "DO NOTHING",
			inTag:                testTag,
			ExpectedError:        nil,
		},
		{
			testName:             "TagAlreadyAssigned",
			inUserID:             1,
			inNoteID:             1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "ASSIGN",
			inTag:                testTag,
			ExpectedError:        nil,
		},
		{
			testName:             "NoteDoesNotExist",
			inUserID:             1,
			inNoteID:             1,
			prepareAccountAction: "DO NOTHING",
			prepareNoteAction:    "DO NOTHING",
			prepareTagAction:     "DO NOTHING",
			inTag:                testTag,
			ExpectedError:        e.ClientNoteError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			nr := repository.NewNoteRepositoryImpl(client, logger)
			ar := repository.NewAccountRepositoryImpl(client, logger)
			tr := repository.NewTagRepositoryImpl(client, logger)

			err = prepareAccountRepo(testSuite.prepareAccountAction, ar)
			if err != nil {
				t.Fatalf("Can't do account pre-test action: %s", err)
			}

			err = prepareNoteRepo(testSuite.prepareNoteAction, testSuite.inUserID, nr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			err = prepareTagRepo(testSuite.prepareTagAction, testSuite.inUserID, testSuite.inNoteID, tr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			service := tag.NewService(tr, nr, logger)

			err = service.Create(testSuite.inUserID, testSuite.inNoteID, &testSuite.inTag)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestIntegration_TagGetAll(t *testing.T) {

	testTag := mother.TagMother()

	logging.Init()
	logger := logging.GetLogger()

	testSuites := []struct {
		testName             string
		prepareAccountAction string
		prepareNoteAction    string
		prepareTagAction     string
		inUserID             int
		inNoteID             int
		inTag                model.Tag
		ExpectedError        error
	}{
		{
			testName:             "TagsFound",
			inUserID:             1,
			inNoteID:             1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "ASSIGN",
			inTag:                testTag,
			ExpectedError:        nil,
		},
		{
			testName:             "TagsNotFound",
			inUserID:             1,
			inNoteID:             1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "DO NOTHING",
			inTag:                testTag,
			ExpectedError:        nil,
		},
		{
			testName:             "NoNotesFound",
			inUserID:             1,
			inNoteID:             1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "DO NOTHING",
			prepareTagAction:     "DO NOTHING",
			inTag:                testTag,
			ExpectedError:        nil,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			nr := repository.NewNoteRepositoryImpl(client, logger)
			ar := repository.NewAccountRepositoryImpl(client, logger)
			tr := repository.NewTagRepositoryImpl(client, logger)

			err = prepareAccountRepo(testSuite.prepareAccountAction, ar)
			if err != nil {
				t.Fatalf("Can't do account pre-test action: %s", err)
			}

			err = prepareNoteRepo(testSuite.prepareNoteAction, testSuite.inUserID, nr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			err = prepareTagRepo(testSuite.prepareTagAction, testSuite.inUserID, testSuite.inNoteID, tr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			service := tag.NewService(tr, nr, logger)

			_, err = service.GetAll(testSuite.inUserID)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestIntegration_TagGetAllByNote(t *testing.T) {

	testTag := mother.TagMother()

	logging.Init()
	logger := logging.GetLogger()

	testSuites := []struct {
		testName             string
		prepareAccountAction string
		prepareNoteAction    string
		prepareTagAction     string
		inUserID             int
		inNoteID             int
		inTag                model.Tag
		ExpectedError        error
	}{
		{
			testName:             "TagsFound",
			inUserID:             1,
			inNoteID:             1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "ASSIGN",
			inTag:                testTag,
			ExpectedError:        nil,
		},
		{
			testName:             "TagsNotFound",
			inUserID:             1,
			inNoteID:             1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "DO NOTHING",
			inTag:                testTag,
			ExpectedError:        nil,
		},
		{
			testName:             "NoteNotFound",
			inUserID:             1,
			inNoteID:             1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "DO NOTHING",
			prepareTagAction:     "DO NOTHING",
			inTag:                testTag,
			ExpectedError:        e.ClientNoteError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			nr := repository.NewNoteRepositoryImpl(client, logger)
			ar := repository.NewAccountRepositoryImpl(client, logger)
			tr := repository.NewTagRepositoryImpl(client, logger)

			err = prepareAccountRepo(testSuite.prepareAccountAction, ar)
			if err != nil {
				t.Fatalf("Can't do account pre-test action: %s", err)
			}

			err = prepareNoteRepo(testSuite.prepareNoteAction, testSuite.inUserID, nr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			err = prepareTagRepo(testSuite.prepareTagAction, testSuite.inUserID, testSuite.inNoteID, tr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			service := tag.NewService(tr, nr, logger)

			_, err = service.GetAllByNote(testSuite.inUserID, testSuite.inNoteID)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestIntegration_TagGetOne(t *testing.T) {

	testTag := mother.TagMother()

	logging.Init()
	logger := logging.GetLogger()

	testSuites := []struct {
		testName             string
		prepareAccountAction string
		prepareNoteAction    string
		prepareTagAction     string
		inUserID             int
		inNoteID             int
		inTagID              int
		inTag                model.Tag
		ExpectedError        error
	}{
		{
			testName:             "TagFound",
			inUserID:             1,
			inNoteID:             1,
			inTagID:              1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "ASSIGN",
			inTag:                testTag,
			ExpectedError:        nil,
		},
		{
			testName:             "TagNotFound",
			inUserID:             1,
			inNoteID:             1,
			inTagID:              1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "DO NOTHING",
			inTag:                testTag,
			ExpectedError:        e.ClientTagError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			nr := repository.NewNoteRepositoryImpl(client, logger)
			ar := repository.NewAccountRepositoryImpl(client, logger)
			tr := repository.NewTagRepositoryImpl(client, logger)

			err = prepareAccountRepo(testSuite.prepareAccountAction, ar)
			if err != nil {
				t.Fatalf("Can't do account pre-test action: %s", err)
			}

			err = prepareNoteRepo(testSuite.prepareNoteAction, testSuite.inUserID, nr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			err = prepareTagRepo(testSuite.prepareTagAction, testSuite.inUserID, testSuite.inNoteID, tr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			service := tag.NewService(tr, nr, logger)

			_, err = service.GetOne(testSuite.inUserID, testSuite.inTagID)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestIntegration_TagUpdate(t *testing.T) {

	testTag := mother.TagMother()

	newTag := testTag
	newTag.Label = "new"

	logging.Init()
	logger := logging.GetLogger()

	testSuites := []struct {
		testName             string
		prepareAccountAction string
		prepareNoteAction    string
		prepareTagAction     string
		inUserID             int
		inNoteID             int
		inTagID              int
		inTag                model.Tag
		ExpectedError        error
	}{
		{
			testName:             "TagFound",
			inUserID:             1,
			inNoteID:             1,
			inTagID:              1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "ASSIGN",
			inTag:                testTag,
			ExpectedError:        nil,
		},
		{
			testName:             "TagNotFound",
			inUserID:             1,
			inNoteID:             1,
			inTagID:              1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "DO NOTHING",
			inTag:                testTag,
			ExpectedError:        e.ClientTagError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			nr := repository.NewNoteRepositoryImpl(client, logger)
			ar := repository.NewAccountRepositoryImpl(client, logger)
			tr := repository.NewTagRepositoryImpl(client, logger)

			err = prepareAccountRepo(testSuite.prepareAccountAction, ar)
			if err != nil {
				t.Fatalf("Can't do account pre-test action: %s", err)
			}

			err = prepareNoteRepo(testSuite.prepareNoteAction, testSuite.inUserID, nr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			err = prepareTagRepo(testSuite.prepareTagAction, testSuite.inUserID, testSuite.inNoteID, tr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			service := tag.NewService(tr, nr, logger)

			err = service.Update(testSuite.inUserID, testSuite.inTagID, newTag)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestIntegration_TagDetach(t *testing.T) {

	testTag := mother.TagMother()

	newTag := testTag
	newTag.Label = "new"

	logging.Init()
	logger := logging.GetLogger()

	testSuites := []struct {
		testName             string
		prepareAccountAction string
		prepareNoteAction    string
		prepareTagAction     string
		inUserID             int
		inAssignedNoteID     int
		inNoteID             int
		inTagID              int
		ExpectedError        error
	}{
		{
			testName:             "TagAttachedToFittingNote",
			inUserID:             1,
			inNoteID:             1,
			inAssignedNoteID:     1,
			inTagID:              1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "ASSIGN",
			ExpectedError:        nil,
		},
		{
			testName:             "TagAttachedToNonFittingNote",
			inUserID:             1,
			inNoteID:             1,
			inAssignedNoteID:     2,
			inTagID:              1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT x2",
			prepareTagAction:     "ASSIGN",
			ExpectedError:        nil,
		},
		{
			testName:             "NoteNotFound",
			inUserID:             1,
			inNoteID:             0,
			inAssignedNoteID:     1,
			inTagID:              1,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "ASSIGN",
			ExpectedError:        e.ClientNoteError,
		},
		{
			testName:             "TagNotFound",
			inUserID:             1,
			inNoteID:             1,
			inAssignedNoteID:     1,
			inTagID:              0,
			prepareAccountAction: "INSERT",
			prepareNoteAction:    "INSERT",
			prepareTagAction:     "ASSIGN",
			ExpectedError:        e.ClientTagError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			client, err := integration.GetTestResource()
			if err != nil {
				t.Fatal(err)
			}

			nr := repository.NewNoteRepositoryImpl(client, logger)
			ar := repository.NewAccountRepositoryImpl(client, logger)
			tr := repository.NewTagRepositoryImpl(client, logger)

			err = prepareAccountRepo(testSuite.prepareAccountAction, ar)
			if err != nil {
				t.Fatalf("Can't do account pre-test action: %s", err)
			}

			err = prepareNoteRepo(testSuite.prepareNoteAction, testSuite.inUserID, nr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			err = prepareTagRepo(testSuite.prepareTagAction, testSuite.inUserID, testSuite.inAssignedNoteID, tr)
			if err != nil {
				t.Fatalf("Can't do pre-test note action: %s", err)
			}

			service := tag.NewService(tr, nr, logger)

			err = service.Detach(testSuite.inUserID, testSuite.inTagID, testSuite.inNoteID)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}
