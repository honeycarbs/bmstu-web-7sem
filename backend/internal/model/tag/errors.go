package tag

type CanNotCreateTagErr struct{}

func (a *CanNotCreateTagErr) Error() string {
	return "can't create tag"
}

type TagNotFoundErr struct{}

func (a *TagNotFoundErr) Error() string {
	return "tag does not exist or does not belong to user"
}
