package user

type RegisterUserDTO struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
