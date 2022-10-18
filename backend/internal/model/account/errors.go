package account

type CanNotCreateAccountErr struct{}

func (a *CanNotCreateAccountErr) Error() string {
	return "can't create accounts"
}

type AccountNotFoundErr struct{}

func (a *AccountNotFoundErr) Error() string {
	return "username or password is incorrect"
}

type CanNotLoginErr struct{}

type CanNotGetErr struct{}

func (a *CanNotGetErr) Error() string {
	return "can't get account"
}

func (a *CanNotLoginErr) Error() string {
	return "can't login"
}

type PasswordDoesNotMatchErr struct{}

func (a *PasswordDoesNotMatchErr) Error() string {
	return "password does not match"
}
