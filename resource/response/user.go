package response

import "time"

type UserLogin struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResetPasswordRequest struct {
	Duration time.Duration `json:"duration"`
	Expire   time.Time     `json:"expire"`
}
