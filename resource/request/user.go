package request

type UserCreate struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
