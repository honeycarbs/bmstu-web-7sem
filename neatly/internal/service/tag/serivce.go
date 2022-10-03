package tag

import (
	"errors"
	"neatly/internal/model/tag"
	"neatly/internal/repository"
	"neatly/pkg/logging"
	"strings"
)

type Service struct {
	tagsRepository  repository.Tag
	notesRepository repository.Note
	logger          logging.Logger
}

func NewService(tr repository.Tag, nr repository.Note, l logging.Logger) *Service {
	return &Service{tagsRepository: tr, notesRepository: nr, logger: l}
}

func (s *Service) Create(userID, noteID int, t *tag.Tag) error {
	_, err := s.notesRepository.GetOne(userID, noteID)
	if err != nil {
		return errors.New("note does not exists or does not belong to auth")
	}

	tags, err := s.tagsRepository.GetAll(userID)
	if err != nil {
		return err
	}

	unique, tuID := s.checkIfUnique(tags, *t)
	if !unique {
		s.logger.Infof("Tag with ID %v is not unique", tuID)
		assigned, err := s.checkIfAssigned(tuID, noteID, userID)
		if err != nil {
			return err
		}
		if !assigned {
			s.logger.Infof("Tag with ID %v is not assigned to note %v", tuID, noteID)
			t.ID = tuID
			err := s.tagsRepository.Assign(tuID, noteID, userID)
			return err
		}
		t.ID = tuID
		return nil
	}

	s.logger.Infof("Tag with ID %v is inuque and will be assigned to note with ID %v", t.ID, noteID)
	err = s.tagsRepository.Create(userID, noteID, t)
	if err != nil {
		return err
	}
	err = s.tagsRepository.Assign(t.ID, noteID, userID)
	return err
}

func (s *Service) GetAll(userID int) ([]tag.Tag, error) {
	return s.tagsRepository.GetAll(userID)
}

func (s *Service) GetAllByNote(userID, noteID int) ([]tag.Tag, error) {
	return s.tagsRepository.GetAllByNote(userID, noteID)
}

func (s *Service) GetOne(userID, tagID int) (tag.Tag, error) {
	return s.tagsRepository.GetOne(userID, tagID)
}

func (s *Service) Delete(userID, tagID int) error {
	return s.tagsRepository.Delete(userID, tagID)
}

func (s *Service) Update(userID, tagID int, t tag.Tag) error {
	tp, err := s.tagsRepository.GetOne(userID, tagID)
	if err != nil {
		return err
	}

	if t.Color == "" {
		t.Color = tp.Color
	}
	if t.Name == "" {
		t.Name = tp.Name
	}

	return s.tagsRepository.Update(userID, tagID, t)
}

func (s *Service) checkIfUnique(tags []tag.Tag, tu tag.Tag) (bool, int) {
	for _, t := range tags {
		if strings.Compare(t.Name, tu.Name) == 0 {
			s.logger.Infof("Found matching tag with id %v", t.ID)
			return false, t.ID
		}
	}
	return true, 0
}

func (s *Service) checkIfAssigned(tagID, noteID, userID int) (bool, error) {
	tags, err := s.tagsRepository.GetAllByNote(userID, noteID)
	if err != nil {
		return false, err
	}

	for _, t := range tags {
		if t.ID == tagID {
			s.logger.Infof("Found matching tag %v assigned to note %v", tagID, noteID)
			return true, nil
		}
	}
	return false, nil
}
