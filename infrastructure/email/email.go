package email

import (
	_ "embed"
	"strconv"

	"gopkg.in/gomail.v2"

	"go-gin-ddd/config"
)

type Body interface {
	Subject() string
	Html() (string, error)
	Plain() (string, error)
}

type IEmail interface {
	Send(to string, body Body) error
}

type email struct{}

func New() IEmail {
	return &email{}
}

func (e email) Send(to string, body Body) error {
	m := gomail.NewMessage()

	html, err := body.Html()
	if err != nil {
		return err
	}
	m.SetBody("text/html", html)

	plain, err := body.Plain()
	if err != nil {
		return err
	}
	m.AddAlternative("text/plain", plain)

	m.SetHeaders(
		map[string][]string{
			"From":    {m.FormatAddress(config.Env.Mail.From, config.Env.Mail.Name)},
			"To":      {to},
			"Subject": {body.Subject()},
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
