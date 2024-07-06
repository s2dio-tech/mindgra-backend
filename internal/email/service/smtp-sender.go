package service

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

type SMTPSender struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (s *SMTPSender) SendEmail(email domain.Email) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	to := common.StreamOf(email.To).Map(
		func(to domain.Contact) string {
			return to.Email
		},
	).Out().([]string)
	msg := []byte(fmt.Sprintf(
		"From: %s<%s>\nTo: %s\nSubject: %s\nmime := \"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n%s\n",
		email.From.Name,
		email.From.Email,
		strings.Join(to, ","),
		email.Template.Subject,
		email.Template.Html,
	))

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", s.Host, s.Port),
		auth,
		email.From.Email,
		to,
		msg,
	)
}
