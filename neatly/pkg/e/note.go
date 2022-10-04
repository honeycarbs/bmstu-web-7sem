package e

type CanNotCreateNoteErr struct{}

func (a *CanNotCreateNoteErr) Error() string {
	return "can't create note"
}

type CanNotAssignNoteErr struct{}

func (a *CanNotAssignNoteErr) Error() string {
	return "can't assign note"
}

type NoteNotFoundErr struct{}

func (a *NoteNotFoundErr) Error() string {
	return "note does not exist or does not belong to user"
}
