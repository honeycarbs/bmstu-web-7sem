package note

import (
	"database/sql"
	"github.com/go-playground/assert/v2"
	"github.com/go-test/deep"
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
	type noteRepoMockBehaviour func(r *mock.MockNoteRepository, n *model.Note)
	testNote := mother.NoteMother()

	testSuites := []struct {
		testName            string
		inNote              model.Note
		CreateNoteBehaviour noteRepoMockBehaviour
		outNote             model.Note
		ExpectedError       error
	}{
		{
			testName: "CreateNoteSuccessful",
			inNote:   testNote,
			CreateNoteBehaviour: func(r *mock.MockNoteRepository, n *model.Note) {
				r.EXPECT().Create(0, n).Return(nil)
			},
			outNote:       testNote,
			ExpectedError: nil,
		},
		{
			testName: "CreateNoteFailure",
			inNote:   testNote,
			CreateNoteBehaviour: func(r *mock.MockNoteRepository, n *model.Note) {
				r.EXPECT().Create(0, n).Return(sql.ErrTxDone)
			},
			outNote:       testNote,
			ExpectedError: sql.ErrTxDone,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoMock := mock.NewMockNoteRepository(c)
			testSuite.CreateNoteBehaviour(repoMock, &testSuite.inNote)

			logging.Init()
			repo := &repository.NoteRepositoryImpl{
				NoteRepository: repoMock,
			}
			mockService := NewService(repo, nil, logging.GetLogger())

			err := mockService.Create(0, &testSuite.inNote)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestService_GetAll(t *testing.T) {
	type noteRepoMockBehaviour func(r *mock.MockNoteRepository, UserID int)
	type tagRepoMockBehaviour func(r *mock.MockTagRepository, UserID, NoteID int)

	testNote := mother.NoteMother()
	testTag := mother.TagMother()

	testNoteWithoutTags := testNote
	testNoteWithoutTags.Tags = []model.Tag{}

	testNoteWithTag := testNote
	testNoteWithTag.Tags = []model.Tag{testTag}

	testSuites := []struct {
		testName          string
		GetNotesBehaviour noteRepoMockBehaviour
		GetTagsBehaviour  tagRepoMockBehaviour
		outNotes          []model.Note
		ExpectedError     error
	}{
		{
			testName: "UserHasNoNotes",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID int) {
				r.EXPECT().GetAll(UserID).Return([]model.Note{}, nil)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Times(0)
			},
			outNotes:      []model.Note{},
			ExpectedError: nil,
		},
		{
			testName: "UserHasNoteWithNoTags",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID int) {
				r.EXPECT().GetAll(UserID).Return([]model.Note{testNote}, nil)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Return([]model.Tag{}, nil)
			},
			outNotes:      []model.Note{testNoteWithoutTags},
			ExpectedError: nil,
		},
		{
			testName: "UserHasNoteWithTag",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID int) {
				r.EXPECT().GetAll(UserID).Return([]model.Note{testNote}, nil)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Return([]model.Tag{testTag}, nil)
			},
			outNotes:      []model.Note{testNoteWithTag},
			ExpectedError: nil,
		},
		{
			testName: "GetNoteError",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID int) {
				r.EXPECT().GetAll(0).Return([]model.Note{}, sql.ErrNoRows)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Times(0)
			},
			outNotes:      []model.Note{},
			ExpectedError: sql.ErrNoRows,
		},
		{
			testName: "GetTagError",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID int) {
				r.EXPECT().GetAll(0).Return([]model.Note{testNote}, nil)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Return([]model.Tag{}, sql.ErrNoRows)
			},
			outNotes:      []model.Note{},
			ExpectedError: sql.ErrNoRows,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noteRepoMock := mock.NewMockNoteRepository(c)
			testSuite.GetNotesBehaviour(noteRepoMock, 0)

			tagRepoMock := mock.NewMockTagRepository(c)
			testSuite.GetTagsBehaviour(tagRepoMock, 0, 0)

			logging.Init()
			noteRepo := &repository.NoteRepositoryImpl{
				NoteRepository: noteRepoMock,
			}

			tagRepo := &repository.TagRepositoryImpl{
				TagRepository: tagRepoMock,
			}
			mockService := NewService(noteRepo, tagRepo, logging.GetLogger())

			got, err := mockService.GetAll(0)

			assert.Equal(t, testSuite.ExpectedError, err)
			if diff := deep.Equal(testSuite.outNotes, got); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestService_GetOne(t *testing.T) {
	type noteRepoMockBehaviour func(r *mock.MockNoteRepository, UserID, noteID int)
	type tagRepoMockBehaviour func(r *mock.MockTagRepository, UserID, NoteID int)

	testNote := mother.NoteMother()
	testTag := mother.TagMother()

	testNoteWithoutTags := testNote
	testNoteWithoutTags.Tags = []model.Tag{}

	testNoteWithTag := testNote
	testNoteWithTag.Tags = []model.Tag{testTag}

	testSuites := []struct {
		testName          string
		GetNotesBehaviour noteRepoMockBehaviour
		GetTagsBehaviour  tagRepoMockBehaviour
		outNote           model.Note
		ExpectedError     error
	}{
		{
			testName: "UserHasNoNotes",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID, noteID int) {
				r.EXPECT().GetOne(UserID, noteID).Return(model.Note{}, e.ClientNoteError)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Times(0)
			},
			outNote:       model.Note{},
			ExpectedError: e.ClientNoteError,
		},
		{
			testName: "UserHasNoteWithNoTags",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID, noteID int) {
				r.EXPECT().GetOne(UserID, noteID).Return(testNote, nil)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Return([]model.Tag{}, nil)
			},
			outNote:       testNoteWithoutTags,
			ExpectedError: nil,
		},
		{
			testName: "UserHasNoteWithTag",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID, noteID int) {
				r.EXPECT().GetOne(UserID, noteID).Return(testNote, nil)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Return([]model.Tag{testTag}, nil)
			},
			outNote:       testNoteWithTag,
			ExpectedError: nil,
		},
		{
			testName: "GetTagError",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID, noteID int) {
				r.EXPECT().GetOne(UserID, noteID).Return(testNote, nil)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Return([]model.Tag{}, sql.ErrNoRows)
			},
			outNote:       model.Note{},
			ExpectedError: sql.ErrNoRows,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noteRepoMock := mock.NewMockNoteRepository(c)
			testSuite.GetNotesBehaviour(noteRepoMock, 0, 0)

			tagRepoMock := mock.NewMockTagRepository(c)
			testSuite.GetTagsBehaviour(tagRepoMock, 0, 0)

			logging.Init()
			noteRepo := &repository.NoteRepositoryImpl{
				NoteRepository: noteRepoMock,
			}

			tagRepo := &repository.TagRepositoryImpl{
				TagRepository: tagRepoMock,
			}
			mockService := NewService(noteRepo, tagRepo, logging.GetLogger())

			got, err := mockService.GetOne(0, 0)

			assert.Equal(t, testSuite.ExpectedError, err)
			if diff := deep.Equal(testSuite.outNote, got); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	type noteRepoMockBehaviour func(r *mock.MockNoteRepository, UserID, noteID int)

	testNote := mother.NoteMother()

	testSuites := []struct {
		testName            string
		DeleteNoteBehaviour noteRepoMockBehaviour
		ExpectedError       error
	}{
		{
			testName: "DeletedSuccessfully",
			DeleteNoteBehaviour: func(r *mock.MockNoteRepository, UserID, noteID int) {
				r.EXPECT().GetOne(UserID, noteID).Return(testNote, nil)
				r.EXPECT().Delete(UserID, noteID).Return(nil)
			},
			ExpectedError: nil,
		},
		{
			testName: "NoteNotFound",
			DeleteNoteBehaviour: func(r *mock.MockNoteRepository, UserID, noteID int) {
				r.EXPECT().GetOne(UserID, noteID).Return(model.Note{}, sql.ErrNoRows)
				r.EXPECT().Delete(UserID, noteID).Times(0)
			},
			ExpectedError: e.ClientNoteError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoMock := mock.NewMockNoteRepository(c)
			testSuite.DeleteNoteBehaviour(repoMock, 0, 0)

			logging.Init()
			repo := &repository.NoteRepositoryImpl{
				NoteRepository: repoMock,
			}
			mockService := NewService(repo, nil, logging.GetLogger())

			err := mockService.Delete(0, 0)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}

func TestService_FindByTags(t *testing.T) {
	type noteRepoMockBehaviour func(r *mock.MockNoteRepository, UserID int)
	type tagRepoMockBehaviour func(r *mock.MockTagRepository, UserID, NoteID int)

	testNote := mother.NoteMother()

	testTag := mother.TagMother()
	testTag.Name = "test"

	testNoteWithoutTags := testNote
	testNoteWithoutTags.Tags = []model.Tag{}

	testNoteWithTag := testNote
	testNoteWithTag.Tags = []model.Tag{testTag}

	testSuites := []struct {
		testName          string
		GetNotesBehaviour noteRepoMockBehaviour
		GetTagsBehaviour  tagRepoMockBehaviour
		outNotes          []model.Note
		ExpectedError     error
	}{
		{
			testName: "UserHasNoNotes",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID int) {
				r.EXPECT().GetAll(UserID).Return([]model.Note{}, nil)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Times(0)
			},
			outNotes:      []model.Note{},
			ExpectedError: nil,
		},
		{
			testName: "UserHasNoteWithNoTags",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID int) {
				r.EXPECT().GetAll(UserID).Return([]model.Note{testNote}, nil)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Return([]model.Tag{}, nil)
			},
			outNotes:      []model.Note{},
			ExpectedError: nil,
		},
		{
			testName: "UserHasNoteWithTagThatSuits",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID int) {
				r.EXPECT().GetAll(UserID).Return([]model.Note{testNote}, nil)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Return([]model.Tag{testTag}, nil)
			},
			outNotes:      []model.Note{testNoteWithTag},
			ExpectedError: nil,
		},
		{
			testName: "GetNoteError",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID int) {
				r.EXPECT().GetAll(0).Return([]model.Note{}, sql.ErrNoRows)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Times(0)
			},
			outNotes:      []model.Note{},
			ExpectedError: sql.ErrNoRows,
		},
		{
			testName: "GetTagError",
			GetNotesBehaviour: func(r *mock.MockNoteRepository, UserID int) {
				r.EXPECT().GetAll(0).Return([]model.Note{testNote}, nil)
			},
			GetTagsBehaviour: func(r *mock.MockTagRepository, UserID, NoteID int) {
				r.EXPECT().GetAllByNote(UserID, NoteID).Return([]model.Tag{}, sql.ErrNoRows)
			},
			outNotes:      []model.Note{},
			ExpectedError: sql.ErrNoRows,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noteRepoMock := mock.NewMockNoteRepository(c)
			testSuite.GetNotesBehaviour(noteRepoMock, 0)

			tagRepoMock := mock.NewMockTagRepository(c)
			testSuite.GetTagsBehaviour(tagRepoMock, 0, 0)

			logging.Init()
			noteRepo := &repository.NoteRepositoryImpl{
				NoteRepository: noteRepoMock,
			}

			tagRepo := &repository.TagRepositoryImpl{
				TagRepository: tagRepoMock,
			}
			mockService := NewService(noteRepo, tagRepo, logging.GetLogger())

			tags := []string{"test"}
			got, err := mockService.FindByTags(0, tags)

			assert.Equal(t, testSuite.ExpectedError, err)
			if diff := deep.Equal(testSuite.outNotes, got); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	type noteRepoMockBehaviour func(r *mock.MockNoteRepository, UserID, noteID int, testNote model.Note)

	testNote := mother.NoteMother()
	noteForUpdate := testNote
	noteForUpdate.Header = "UPDATE"

	testSuites := []struct {
		testName            string
		UpdateNoteBehaviour noteRepoMockBehaviour
		inNote              model.Note
		outNote             model.Note
		needsBodyUpdate     bool
		ExpectedError       error
	}{
		{
			testName: "NoteUpdatedWithoutBody",
			UpdateNoteBehaviour: func(r *mock.MockNoteRepository, UserID, noteID int, n model.Note) {
				r.EXPECT().GetOne(UserID, noteID).Return(noteForUpdate, nil)
				r.EXPECT().Update(UserID, noteForUpdate).Return(nil)
			},
			needsBodyUpdate: false,
			ExpectedError:   nil,
		},
		{
			testName: "NoteUpdatedWithBody",
			UpdateNoteBehaviour: func(r *mock.MockNoteRepository, UserID, noteID int, n model.Note) {
				r.EXPECT().GetOne(UserID, noteID).Return(noteForUpdate, nil)
				r.EXPECT().Update(UserID, noteForUpdate).Return(nil)
			},
			needsBodyUpdate: true,
			ExpectedError:   nil,
		},
		{
			testName: "NoteDoesNotExist",
			UpdateNoteBehaviour: func(r *mock.MockNoteRepository, UserID, noteID int, n model.Note) {
				r.EXPECT().GetOne(UserID, noteID).Return(model.Note{}, e.ClientNoteError)
				r.EXPECT().Update(UserID, noteForUpdate).Times(0)
			},
			needsBodyUpdate: true,
			ExpectedError:   e.ClientNoteError,
		},
	}
	for _, testSuite := range testSuites {
		t.Run(testSuite.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoMock := mock.NewMockNoteRepository(c)
			testSuite.UpdateNoteBehaviour(repoMock, 0, 0, testNote)

			logging.Init()
			repo := &repository.NoteRepositoryImpl{
				NoteRepository: repoMock,
			}
			mockService := NewService(repo, nil, logging.GetLogger())

			err := mockService.Update(0, testSuite.inNote, testSuite.needsBodyUpdate)

			assert.Equal(t, testSuite.ExpectedError, err)
		})
	}
}
