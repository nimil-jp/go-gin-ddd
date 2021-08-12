package repository

type IEmail interface {
	Send(to string, subject string, body string) error
}
