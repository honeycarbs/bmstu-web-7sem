package tag

import (
	"errors"
	"neatly/internal/model"
	"neatly/internal/repository"
	"neatly/pkg/logging"
	"strings"
)

type Service struct {
	tagsRepository  *repository.TagRepositoryImpl
	notesRepository *repository.NoteRepositoryImpl
	logger          logging.Logger
}

func NewService(tagsRepository *repository.TagRepositoryImpl, notesRepository *repository.NoteRepositoryImpl, logger logging.Logger) *Service {
	return &Service{tagsRepository: tagsRepository, notesRepository: notesRepository, logger: logger}
}

func (s *Service) Create(userID, noteID int, t *model.Tag) error {
	_, err := s.notesRepository.GetOne(userID, noteID)
	if err != nil {
		return errors.New("note does not exist or does not belong to accounts")
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

func (s *Service) GetAll(userID int) ([]model.Tag, error) {
	return s.tagsRepository.GetAll(userID)
}

func (s *Service) GetAllByNote(userID, noteID int) ([]model.Tag, error) {
	return s.tagsRepository.GetAllByNote(userID, noteID)
}

func (s *Service) GetOne(userID, tagID int) (model.Tag, error) {
	return s.tagsRepository.GetOne(userID, tagID)
}

func (s *Service) Delete(userID, tagID int) error {
	return s.tagsRepository.Delete(userID, tagID)
}

func (s *Service) Update(userID, tagID int, t model.Tag) error {
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

func (s *Service) Detach(userID, tagID, noteID int) error {
	_, err := s.notesRepository.GetOne(userID, noteID)
	if err != nil {
		return errors.New("note does not exist or does not belong to user")
	}

	ns, err := s.notesRepository.GetAll(userID)
	if err != nil {
		return err
	}

	t, err := s.tagsRepository.GetOne(userID, tagID)
	s.logger.Infof("Found tag %v: %v, %v", tagID, t.Name, t.Color)
	if err != nil {
		return err
	}

	for _, n := range ns {
		n.Tags, err = s.tagsRepository.GetAllByNote(userID, n.ID)
		if err != nil {
			return err
		}
		if n.HasSpecificTag(t.Name) && n.ID != noteID {
			s.logger.Infof("Found this tag at note %v", n.ID)
			err = s.tagsRepository.Detach(userID, tagID, noteID)
			return err
		}
	}

	s.logger.Info("Deleting tag")
	err = s.tagsRepository.Delete(userID, tagID)
	return err
}

func (s *Service) checkIfUnique(tags []model.Tag, tu model.Tag) (bool, int) {
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
