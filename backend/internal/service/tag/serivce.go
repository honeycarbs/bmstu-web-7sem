package tag

import (
	"neatly/internal/model"
	"neatly/internal/repository"
	"neatly/pkg/e"
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

func (s *Service) Create(userID, noteID int, t *model.Tag) (bool, error) {
	_, err := s.notesRepository.GetOne(userID, noteID)
	if err != nil {
		return false, e.ClientNoteError
	}

	tags, err := s.tagsRepository.GetAll(userID)
	if err != nil {
		return false, err
	}

	modified := false
	unique, tuID := s.checkIfUnique(tags, *t)
	if !unique {
		s.logger.Infof("Tag with ID %v is not unique", tuID)
		assigned, err := s.checkIfAssigned(tuID, noteID, userID)
		if err != nil {
			return modified, err
		}
		if !assigned {
			modified = true
			s.logger.Infof("Tag with ID %v is not assigned to note %v", tuID, noteID)
			t.ID = tuID
			err := s.tagsRepository.Assign(tuID, noteID, userID)
			return modified, err
		}
		t.ID = tuID
		return modified, nil
	}

	s.logger.Infof("Tag with ID %v is inuque and will be assigned to note with ID %v", t.ID, noteID)
	modified = true
	err = s.tagsRepository.Create(userID, noteID, t)
	if err != nil {
		return false, err
	}
	err = s.tagsRepository.Assign(t.ID, noteID, userID)
	return modified, err
}

func (s *Service) GetAll(userID int) ([]model.Tag, error) {
	return s.tagsRepository.GetAll(userID)
}

func (s *Service) GetAllByNote(userID, noteID int) ([]model.Tag, error) {
	_, err := s.notesRepository.GetOne(userID, noteID)

	if err != nil {
		return []model.Tag{}, e.ClientNoteError
	}

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
		return e.ClientTagError
	}

	if t.Label == "" {
		t.Label = tp.Label
	}

	return s.tagsRepository.Update(userID, tagID, t)
}

func (s *Service) Detach(userID, tagID, noteID int) error {
	inNote, err := s.notesRepository.GetOne(userID, noteID)
	if err != nil {
		return e.ClientNoteError
	}
	inTag, err := s.tagsRepository.GetOne(userID, tagID)
	if err != nil {
		return e.ClientTagError
	}

	inNote.Tags, err = s.tagsRepository.GetAllByNote(userID, noteID)
	if err != nil {
		return err
	}

	var (
		attachedToMany = false
	)

	if inNote.HasSpecificTag(inTag.Label) {
		s.logger.Info("Detaching tag...")
		err = s.tagsRepository.Detach(userID, tagID, noteID)
	} else {
		s.logger.Info("Tag is not attached to this note.")
		return nil
	}

	ns, err := s.notesRepository.GetAll(userID)
	if err != nil {
		return err
	}
	for _, n := range ns {
		if n.ID != noteID {
			n.Tags, err = s.tagsRepository.GetAllByNote(userID, n.ID)
			if n.HasSpecificTag(inTag.Label) {
				attachedToMany = true
			}
		}
	}

	if !attachedToMany {
		s.logger.Info("Tag is attached to one note and should be deleted.")
		err = s.tagsRepository.Delete(userID, tagID)
		return err
	}
	return nil
}

func (s *Service) checkIfUnique(tags []model.Tag, tu model.Tag) (bool, int) {
	for _, t := range tags {
		if strings.Compare(t.Label, tu.Label) == 0 {
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
