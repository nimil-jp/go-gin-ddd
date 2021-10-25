package email

import "fmt"

type UserResetPasswordRequest struct {
	Url   string
	Token string
}

func (r UserResetPasswordRequest) Subject() string {
	return "パスワードリセット"
}

func (r UserResetPasswordRequest) Html() (string, error) {
	return fmt.Sprintf("url: %s\ntoken: %s", r.Url, r.Token), nil
}

func (r UserResetPasswordRequest) Plain() (string, error) {
	return fmt.Sprintf("url: %s\ntoken: %s", r.Url, r.Token), nil
}
