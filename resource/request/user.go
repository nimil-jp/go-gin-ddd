package request

type UserCreate struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UserLogin struct {
	Session  bool   `json:"session"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRefreshToken struct {
	Session      bool   `json:"session"`
	RefreshToken string `json:"refresh_token"`
}

type UserResetPasswordRequest struct {
	Email string `json:"email"`
}

type UserResetPassword struct {
	RecoveryToken   string `json:"recovery_token"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}
