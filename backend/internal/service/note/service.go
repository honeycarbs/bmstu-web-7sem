package note

import (
	"neatly/internal/model/note"
	"neatly/internal/repository"
	"neatly/pkg/logging"
)

type Service struct {
	notesRepository repository.Note
	tagsRepository  repository.Tag
	logger          logging.Logger
}

func NewService(notesRepository repository.Note, tagsRepository repository.Tag, logger logging.Logger) *Service {
	return &Service{notesRepository: notesRepository, tagsRepository: tagsRepository, logger: logger}
}

func (s *Service) Create(userID int, n *note.Note) error {
	err := s.notesRepository.Create(userID, n)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAll(userID int) ([]note.Note, error) {
	notes, err := s.notesRepository.GetAll(userID)
	if err != nil {
		return []note.Note{}, err
	}

	for i := 0; i < len(notes); i++ {
		noteID := notes[i].ID
		tags, err := s.tagsRepository.GetAllByNote(userID, noteID)
		if err != nil {
			return []note.Note{}, err
		}
		notes[i].Tags = tags
	}

	return notes, nil
}

func (s *Service) GetOne(userID, noteID int) (note.Note, error) {
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
	return s.notesRepository.Delete(userID, noteID)
}

func (s *Service) Update(userID int, n note.Note, needBodyUpdate bool) error {
	prev, err := s.notesRepository.GetOne(userID, n.ID)
	if err != nil {
		return err
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

func (s *Service) FindByTags(userID int, tagNames []string) ([]note.Note, error) {
	ns, err := s.notesRepository.GetAll(userID)
	if err != nil {
		return ns, err
	}

	var (
		notesWithAllTags []note.Note
	)

	for _, n := range ns {
		n.Tags, err = s.tagsRepository.GetAllByNote(userID, n.ID)
		if err != nil {
			return ns, err
		}

		s.logger.Infof("Found tags from note %v: %v", n.ID, n.Tags)
		if n.HasEveryTag(tagNames) {
			notesWithAllTags = append(notesWithAllTags, n)
		}
	}

	return notesWithAllTags, nil
}
