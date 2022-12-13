package note

import (
	"neatly/internal/model"
	"neatly/internal/repository"
	"neatly/pkg/e"
	"neatly/pkg/logging"
)

type Service struct {
	notesRepository *repository.NoteRepositoryImpl
	tagsRepository  *repository.TagRepositoryImpl
	logger          logging.Logger
}

func NewService(notesRepository *repository.NoteRepositoryImpl, tagsRepository *repository.TagRepositoryImpl, logger logging.Logger) *Service {
	return &Service{notesRepository: notesRepository, tagsRepository: tagsRepository, logger: logger}
}

func (s *Service) Create(userID int, n *model.Note) error {
	err := s.notesRepository.Create(userID, n)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAll(userID int) ([]model.Note, error) {
	notes, err := s.notesRepository.GetAll(userID)
	if err != nil {
		return []model.Note{}, err
	}

	for i := 0; i < len(notes); i++ {
		noteID := notes[i].ID
		tags, err := s.tagsRepository.GetAllByNote(userID, noteID)
		if err != nil {
			return []model.Note{}, err
		}
		notes[i].Tags = tags
	}

	return notes, nil
}

func (s *Service) GetOne(userID, noteID int) (model.Note, error) {
	n, err := s.notesRepository.GetOne(userID, noteID)
	if err != nil {
		return n, err
	}

	tags, err := s.tagsRepository.GetAllByNote(userID, n.ID)
	if err != nil {
		return n, err
	}
	n.Tags = tags

	return n, nil
}

func (s *Service) Delete(userID, noteID int) error {
	_, err := s.notesRepository.GetOne(userID, noteID)
	if err != nil {
		return e.ClientNoteError
	}
	return s.notesRepository.Delete(noteID, userID)
}

func (s *Service) Update(userID int, n model.Note, needBodyUpdate bool) error {
	prev, err := s.notesRepository.GetOne(userID, n.ID)
	if err != nil {
		return e.ClientNoteError
	}
	if n.Header == "" {
		n.Header = prev.Header
	}

	if !needBodyUpdate {
		n.Body = prev.Body
		n.ShortBody = prev.ShortBody
	}

	return s.notesRepository.Update(userID, n)
}

func (s *Service) FindByTags(userID int, tagNames []string) ([]model.Note, error) {
	ns, err := s.notesRepository.GetAll(userID)
	if err != nil {
		return ns, err
	}

	var (
		notesWithAllTags = make([]model.Note, 0)
	)

	for _, n := range ns {
		n.Tags, err = s.tagsRepository.GetAllByNote(userID, n.ID)
		if err != nil {
			return notesWithAllTags, err
		}

		s.logger.Infof("Found tags from note %v: %v", n.ID, n.Tags)
		if n.HasEveryTag(tagNames) {
			notesWithAllTags = append(notesWithAllTags, n)
		}
	}

	return notesWithAllTags, nil
}
