package account

type RegisterAccountDTO struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginAccountDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type WithTokenDTO struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetAccountDTO struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
