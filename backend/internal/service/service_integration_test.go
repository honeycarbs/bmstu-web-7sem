package service_test

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"neatly/internal/model"
	"neatly/internal/model/mother"
	"neatly/internal/repository"
	"neatly/internal/service/account"
	"neatly/internal/service/note"
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

func prepareNoteRepo(action string, userID int, repo repository.NoteRepository) error {
	switch action {
	case "INSERT":
		{
			n := mother.NoteMother()
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

//func TestGenerateJWT(t *testing.T) {
//	testAccount := mother.AccountMother()
//	testAccountInvalidPassword := testAccount
//	testAccountInvalidPassword.Password = "kto prochital tot loh"
//
//	testSuites := []struct {
//		testName         string
//		inAccount        model.Account
//		outAccount       model.Account
//		prepareAction    string
//		ExpectedError    error
//		ExpectedTokenVal string
//	}{
//		{
//			testName:         "AuthorizeSuccessful",
//			inAccount:        testAccount,
//			outAccount:       testAccount,
//			prepareAction:    "INSERT",
//			ExpectedError:    nil,
//			ExpectedTokenVal: mother.TokenMother(),
//		},
//		{
//			testName:         "PasswordDoesNotMatch",
//			inAccount:        testAccountInvalidPassword,
//			outAccount:       testAccount,
//			prepareAction:    "INSERT",
//			ExpectedError:    errors.New("password does not match"),
//			ExpectedTokenVal: "",
//		},
//	}
//
//	logging.Init()
//	logger := logging.GetLogger()
//
//	for _, testSuite := range testSuites {
//		t.Run(testSuite.testName, func(t *testing.T) {
//
//			client, err := integration.GetTestResource()
//			if err != nil {
//				t.Fatal(err)
//			}
//
//			repo := repository.NewAccountRepositoryImpl(client, logger)
//			err = prepareAccountRepo(testSuite.prepareAction, repo)
//			if err != nil {
//				t.Fatalf("Can't do pre-test action: %s", err)
//			}
//
//			service := account.NewService(repo, logger)
//
//			token, err := service.GenerateJWT(&testSuite.inAccount)
//
//			assert.Equal(t, testSuite.ExpectedError, err)
//			assert.Equal(t, testSuite.ExpectedTokenVal, token)
//		})
//	}
//}

func TestAccountGetOne(t *testing.T) {
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

func TestNoteCreate(t *testing.T) {
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

func TestNoteGetOne(t *testing.T) {
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

func TestNoteFindByTag(t *testing.T) {
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

func TestNoteUpdate(t *testing.T) {
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
