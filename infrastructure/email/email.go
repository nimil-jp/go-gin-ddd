package email

import (
	"strconv"

	"gopkg.in/gomail.v2"

	"go-gin-ddd/config"
)

type IEmail interface {
	Send(to string, subject string, body string) error
}

type email struct{}

func New() IEmail {
	return &email{}
}

func (e email) Send(to string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetBody("text/html", body)

	m.AddAlternative("text/plain", body)

	m.SetHeaders(
		map[string][]string{
			"From":    {m.FormatAddress(config.Env.Mail.From, config.Env.Mail.Name)},
			"To":      {to},
			"Subject": {subject},
		},
	)

	port, err := strconv.Atoi(config.Env.Smtp.Port)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(config.Env.Smtp.Host, port, config.Env.Smtp.User, config.Env.Smtp.Password)

	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
