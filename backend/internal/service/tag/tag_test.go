//go:build unit
// +build unit

package tag

import (
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"neatly/internal/model"
	"neatly/internal/model/mother"
	"neatly/internal/repository"
	"neatly/internal/repository/mock"
	"neatly/pkg/e"
	"neatly/pkg/logging"

	"testing"
)

func TestService_Create(t *testing.T) {
	type noteRepoMockBehaviour func(r *mock.MockNoteRepository, UserID, NoteID int)
	type tagRepoMockBehaviour func(r *mock.MockTagRepository, UserID, NoteID int)

	testTag := mother.TagMother()
	testTagUnique := testTag
	testTagUnique.Label = "unique"

	testNote := mother.NoteMother()

	testSuites := []struct {
		testName           string
		inUserID           int
		inNoteID           int
		notesRepoBehaviour noteRepoMockBehaviour
		tagsRepoBehaviour  tagRepoMockBehaviour
		inTag              model.Tag
		ExpectedError      error
	}{
		{
			testName: "TagIsUniqueAndNotAssigned",
			inUserID: 0,
			inNoteID: 0,
			notesRepoBehaviour: func(r *mock.MockNoteRepository, UserID, NoteID int) {
				r.EXPECT().GetOne(UserID, NoteID).Return(testNote, nil)
			},
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAll(UserID).Return([]model.Tag{testTag}, nil)
				r.EXPECT().GetAllByNote(UserID, NoteID).Times(0)
				r.EXPECT().Assign(0, NoteID, UserID).Return(nil)
				r.EXPECT().Create(UserID, NoteID, &testTagUnique).Return(nil)
			},
			inTag:         testTagUnique,
			ExpectedError: nil,
		},
		{
			testName: "TagIsNotUniqueAndNotAssigned",
			inUserID: 0,
			inNoteID: 0,
			notesRepoBehaviour: func(r *mock.MockNoteRepository, UserID, NoteID int) {
				r.EXPECT().GetOne(UserID, NoteID).Return(testNote, nil)
			},
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAll(UserID).Return([]model.Tag{testTag}, nil)
				r.EXPECT().GetAllByNote(UserID, NoteID).Return([]model.Tag{}, nil)
				r.EXPECT().Assign(0, NoteID, UserID).Return(nil)
				r.EXPECT().Create(UserID, NoteID, &testTag).Times(0)
			},
			inTag:         testTag,
			ExpectedError: nil,
		},
		{
			testName: "TagIsNotUniqueAndAssigned",
			inUserID: 0,
			inNoteID: 0,
			notesRepoBehaviour: func(r *mock.MockNoteRepository, UserID, NoteID int) {
				r.EXPECT().GetOne(UserID, NoteID).Return(testNote, nil)
			},
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAll(UserID).Return([]model.Tag{testTag}, nil)
				r.EXPECT().GetAllByNote(UserID, NoteID).Return([]model.Tag{testTag}, nil)
				r.EXPECT().Assign(0, NoteID, UserID).Times(0)
				r.EXPECT().Create(UserID, NoteID, &testTag).Times(0)
			},
			inTag:         testTag,
			ExpectedError: nil,
		},
		{
			testName: "NoteNotFound",
			inUserID: 0,
			inNoteID: 0,
			notesRepoBehaviour: func(r *mock.MockNoteRepository, UserID, NoteID int) {
				r.EXPECT().GetOne(UserID, NoteID).Return(model.Note{}, e.ClientNoteError)
			},
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
			},
			inTag:         testTag,
			ExpectedError: e.ClientNoteError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noteRepoMock := mock.NewMockNoteRepository(c)
			testSuite.notesRepoBehaviour(noteRepoMock, 0, 0)

			tagRepoMock := mock.NewMockTagRepository(c)
			testSuite.tagsRepoBehaviour(tagRepoMock, 0, 0)

			logging.Init()
			noteRepo := &repository.NoteRepositoryImpl{
				NoteRepository: noteRepoMock,
			}

			tagRepo := &repository.TagRepositoryImpl{
				TagRepository: tagRepoMock,
			}
			mockService := NewService(tagRepo, noteRepo, logging.GetLogger())

			err := mockService.Create(0, 0, &testSuite.inTag)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestService_GetAll(t *testing.T) {
	type tagRepoMockBehaviour func(r *mock.MockTagRepository, UserID int)

	testTag := mother.TagMother()

	testSuites := []struct {
		testName          string
		inUserID          int
		tagsRepoBehaviour tagRepoMockBehaviour
		ExpectedError     error
	}{
		{
			testName: "FoundAtLeastOneTag",
			inUserID: 0,
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID int) {
				r.EXPECT().GetAll(UserID).Return([]model.Tag{testTag}, nil)
			},
			ExpectedError: nil,
		},
		{
			testName: "FoundNoTags",
			inUserID: 0,
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID int) {
				r.EXPECT().GetAll(UserID).Return([]model.Tag{}, nil)
			},
			ExpectedError: nil,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			tagRepoMock := mock.NewMockTagRepository(c)
			testSuite.tagsRepoBehaviour(tagRepoMock, 0)

			logging.Init()

			tagRepo := &repository.TagRepositoryImpl{
				TagRepository: tagRepoMock,
			}
			mockService := NewService(tagRepo, nil, logging.GetLogger())

			_, err := mockService.GetAll(testSuite.inUserID)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestService_GetAllByNote(t *testing.T) {
	type tagRepoMockBehaviour func(r *mock.MockTagRepository, UserID, noteID int)

	testTag := mother.TagMother()

	testSuites := []struct {
		testName          string
		inUserID          int
		inNoteID          int
		tagsRepoBehaviour tagRepoMockBehaviour
		ExpectedError     error
	}{
		{
			testName: "FoundAtLeastOneTag",
			inUserID: 0,
			inNoteID: 0,
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, noteID int) {
				r.EXPECT().GetAllByNote(UserID, noteID).Return([]model.Tag{testTag}, nil)
			},
			ExpectedError: nil,
		},
		{
			testName: "FoundNoTags",
			inUserID: 0,
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, noteID int) {
				r.EXPECT().GetAllByNote(UserID, noteID).Return([]model.Tag{}, nil)
			},
			ExpectedError: nil,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			tagRepoMock := mock.NewMockTagRepository(c)
			testSuite.tagsRepoBehaviour(tagRepoMock, 0, 0)

			logging.Init()

			tagRepo := &repository.TagRepositoryImpl{
				TagRepository: tagRepoMock,
			}
			mockService := NewService(tagRepo, nil, logging.GetLogger())

			_, err := mockService.GetAllByNote(testSuite.inUserID, testSuite.inNoteID)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestService_GetOne(t *testing.T) {
	type tagRepoMockBehaviour func(r *mock.MockTagRepository, UserID, tagID int)

	testTag := mother.TagMother()

	testSuites := []struct {
		testName          string
		inUserID          int
		inTagID           int
		tagsRepoBehaviour tagRepoMockBehaviour
		ExpectedError     error
	}{
		{
			testName: "FoundTag",
			inUserID: 0,
			inTagID:  0,
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, tagID int) {
				r.EXPECT().GetOne(UserID, tagID).Return(testTag, nil)
			},
			ExpectedError: nil,
		},
		{
			testName: "FoundNoTags",
			inUserID: 0,
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, tagID int) {
				r.EXPECT().GetOne(UserID, tagID).Return(model.Tag{}, e.ClientTagError)
			},
			ExpectedError: e.ClientTagError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			tagRepoMock := mock.NewMockTagRepository(c)
			testSuite.tagsRepoBehaviour(tagRepoMock, 0, 0)

			logging.Init()

			tagRepo := &repository.TagRepositoryImpl{
				TagRepository: tagRepoMock,
			}
			mockService := NewService(tagRepo, nil, logging.GetLogger())

			_, err := mockService.GetOne(testSuite.inUserID, testSuite.inTagID)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestService_Update(t *testing.T) {
	type tagRepoMockBehaviour func(r *mock.MockTagRepository, UserID, tagID int)

	testTagBeforeUpdate := mother.TagMother()
	testTagBeforeUpdate.Label = "old name"
	testTagBeforeUpdate.Color = "old color"

	testTagNameUpdate := mother.TagMother()
	testTagNameUpdate.Label = "new name"

	testTagNameUpdateFull := testTagNameUpdate
	testTagNameUpdateFull.Color = testTagBeforeUpdate.Color

	testTagColorUpdate := mother.TagMother()
	testTagColorUpdate.Color = "new color"

	testTagColorUpdateFull := testTagColorUpdate
	testTagColorUpdateFull.Label = testTagBeforeUpdate.Label

	testSuites := []struct {
		testName          string
		inUserID          int
		inTagID           int
		inTag             model.Tag
		tagsRepoBehaviour tagRepoMockBehaviour
		ExpectedError     error
	}{
		{
			testName: "NeedColorUpdate",
			inUserID: 0,
			inTagID:  0,
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, tagID int) {
				r.EXPECT().GetOne(UserID, tagID).Return(testTagBeforeUpdate, nil)
				r.EXPECT().Update(UserID, tagID, testTagColorUpdateFull).Return(nil)
			},
			inTag:         testTagColorUpdate,
			ExpectedError: nil,
		},
		{
			testName: "NeedNameUpdate",
			inUserID: 0,
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, tagID int) {
				r.EXPECT().GetOne(UserID, tagID).Return(testTagBeforeUpdate, nil)
				r.EXPECT().Update(UserID, tagID, testTagNameUpdateFull).Return(nil)
			},
			inTag:         testTagNameUpdate,
			ExpectedError: nil,
		},
		{
			testName: "TagNotFound",
			inUserID: 0,
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, tagID int) {
				r.EXPECT().GetOne(UserID, tagID).Return(model.Tag{}, e.ClientTagError)
			},
			inTag:         testTagNameUpdate,
			ExpectedError: e.ClientTagError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			tagRepoMock := mock.NewMockTagRepository(c)
			testSuite.tagsRepoBehaviour(tagRepoMock, 0, 0)

			logging.Init()

			tagRepo := &repository.TagRepositoryImpl{
				TagRepository: tagRepoMock,
			}
			mockService := NewService(tagRepo, nil, logging.GetLogger())

			err := mockService.Update(testSuite.inUserID, testSuite.inTagID, testSuite.inTag)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestService_Delete(t *testing.T) {
	type tagRepoMockBehaviour func(r *mock.MockTagRepository, UserID, tagID int)

	testSuites := []struct {
		testName          string
		inUserID          int
		inTagID           int
		tagsRepoBehaviour tagRepoMockBehaviour
		ExpectedError     error
	}{
		{
			testName: "DeletedSuccessfully",
			inUserID: 0,
			inTagID:  0,
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, tagID int) {
				r.EXPECT().Delete(UserID, tagID).Return(nil)
			},
			ExpectedError: nil,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			tagRepoMock := mock.NewMockTagRepository(c)
			testSuite.tagsRepoBehaviour(tagRepoMock, 0, 0)

			logging.Init()

			tagRepo := &repository.TagRepositoryImpl{
				TagRepository: tagRepoMock,
			}
			mockService := NewService(tagRepo, nil, logging.GetLogger())

			err := mockService.Delete(testSuite.inUserID, testSuite.inTagID)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestService_Detach(t *testing.T) {
	type noteRepoMockBehaviour func(r *mock.MockNoteRepository, UserID, NoteID int)
	type tagRepoMockBehaviour func(r *mock.MockTagRepository, UserID, NoteID, tagID int)

	testTag := mother.TagMother()
	testTag.Label = "name"

	testNoteNoTag := mother.NoteMother()
	testNoteNoTag.ID = 1

	testNoteWithTag := mother.NoteMother()
	testNoteWithTag.Tags = []model.Tag{testTag}

	otherNoteWithTag := mother.NoteMother()
	otherNoteWithTag.Tags = []model.Tag{testTag}
	otherNoteWithTag.ID = 2

	testNotesSliceOneFits := []model.Note{testNoteNoTag, testNoteWithTag}
	testNotesSliceMultipleFits := []model.Note{testNoteNoTag, testNoteWithTag, otherNoteWithTag}

	testSuites := []struct {
		testName           string
		inUserID           int
		inTagID            int
		inNoteID           int
		notesRepoBehaviour noteRepoMockBehaviour
		tagsRepoBehaviour  tagRepoMockBehaviour
		ExpectedError      error
	}{
		{
			testName: "TagAttachedToFittingNoteAndShouldBeDeleted",
			inUserID: 0,
			inTagID:  0,
			inNoteID: 0,
			notesRepoBehaviour: func(r *mock.MockNoteRepository, UserID, NoteID int) {
				r.EXPECT().GetOne(UserID, NoteID).Return(testNoteWithTag, nil)
				r.EXPECT().GetAll(UserID).Return(testNotesSliceOneFits, nil)
			},
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, NoteID, tagID int) {
				r.EXPECT().GetAllByNote(UserID, testNoteWithTag.ID).Return([]model.Tag{testTag}, nil)
				r.EXPECT().GetAllByNote(UserID, testNoteNoTag.ID).Return([]model.Tag{}, nil)
				r.EXPECT().GetOne(UserID, tagID).Return(testTag, nil)
				r.EXPECT().Detach(UserID, tagID, NoteID).Return(nil)
				r.EXPECT().Delete(UserID, tagID).Return(nil)
			},
			ExpectedError: nil,
		},
		{
			testName: "TagAttachedToFittingNoteAndShouldNotBeDeleted",
			inUserID: 0,
			inTagID:  0,
			inNoteID: 0,
			notesRepoBehaviour: func(r *mock.MockNoteRepository, UserID, NoteID int) {
				r.EXPECT().GetOne(UserID, NoteID).Return(testNoteWithTag, nil)
				r.EXPECT().GetAll(UserID).Return(testNotesSliceMultipleFits, nil)
			},
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, NoteID, tagID int) {
				r.EXPECT().GetAllByNote(UserID, testNoteWithTag.ID).Return([]model.Tag{testTag}, nil)
				r.EXPECT().GetAllByNote(UserID, otherNoteWithTag.ID).Return([]model.Tag{testTag}, nil)
				r.EXPECT().GetAllByNote(UserID, testNoteNoTag.ID).Return([]model.Tag{}, nil)
				r.EXPECT().GetOne(UserID, tagID).Return(testTag, nil)
				r.EXPECT().Detach(UserID, tagID, NoteID).Return(nil)
			},
			ExpectedError: nil,
		},
		{
			testName: "TagAttachedToNonFittingNote",
			inUserID: 0,
			inTagID:  0,
			inNoteID: 1,
			notesRepoBehaviour: func(r *mock.MockNoteRepository, UserID, NoteID int) {
				r.EXPECT().GetOne(UserID, NoteID).Return(testNoteWithTag, nil)
			},
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, NoteID, tagID int) {
				r.EXPECT().GetOne(UserID, tagID).Return(model.Tag{}, nil)
				r.EXPECT().GetAllByNote(UserID, testNoteNoTag.ID).Return([]model.Tag{}, nil)
			},
			ExpectedError: nil,
		},
		{
			testName: "NoteNotFound",
			inUserID: 0,
			inTagID:  0,
			inNoteID: 1,
			notesRepoBehaviour: func(r *mock.MockNoteRepository, UserID, NoteID int) {
				r.EXPECT().GetOne(UserID, NoteID).Return(model.Note{}, e.ClientNoteError)
			},
			tagsRepoBehaviour: func(r *mock.MockTagRepository, UserID, NoteID, tagID int) {

			},
			ExpectedError: e.ClientNoteError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noteRepoMock := mock.NewMockNoteRepository(c)
			testSuite.notesRepoBehaviour(noteRepoMock, testSuite.inUserID, testSuite.inNoteID)

			tagRepoMock := mock.NewMockTagRepository(c)
			testSuite.tagsRepoBehaviour(tagRepoMock, testSuite.inUserID, testSuite.inNoteID, testSuite.inTagID)

			logging.Init()
			noteRepo := &repository.NoteRepositoryImpl{
				NoteRepository: noteRepoMock,
			}

			tagRepo := &repository.TagRepositoryImpl{
				TagRepository: tagRepoMock,
			}
			mockService := NewService(tagRepo, noteRepo, logging.GetLogger())

			err := mockService.Detach(testSuite.inUserID, testSuite.inTagID, testSuite.inNoteID)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}
