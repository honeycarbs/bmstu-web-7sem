package auth

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

type JwtDTO struct {
	Token string `json:"token"`
}
