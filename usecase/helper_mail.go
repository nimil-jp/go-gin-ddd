package usecase

import (
	"strconv"

	"go-ddd/config"
	"gopkg.in/gomail.v2"
)

func sendMail(to string, subject string, body string) error {
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

	port, err := strconv.Atoi(config.Env.SMTP.Port)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(config.Env.SMTP.Host, port, config.Env.SMTP.User, config.Env.SMTP.Password)

	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
