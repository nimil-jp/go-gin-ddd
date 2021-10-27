package email

import "fmt"

type UserResetPasswordRequest struct {
	URL   string
	Token string
}

func (r UserResetPasswordRequest) Subject() string {
	return "パスワードリセット"
}

func (r UserResetPasswordRequest) HTML() (string, error) {
	return fmt.Sprintf("url: %s\ntoken: %s", r.URL, r.Token), nil
}

func (r UserResetPasswordRequest) Plain() (string, error) {
	return fmt.Sprintf("url: %s\ntoken: %s", r.URL, r.Token), nil
}
