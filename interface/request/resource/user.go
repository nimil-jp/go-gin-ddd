package resource

type UserAddRequest struct {
	UserID          string `json:"user_id"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}
