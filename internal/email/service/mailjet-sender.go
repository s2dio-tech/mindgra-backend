package service

import (
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

type MailJet struct {
	PublicKey  string
	PrivateKey string
}

func (m *MailJet) SendEmail(email domain.Email) error {
	client := mailjet.NewMailjetClient(m.PublicKey, m.PrivateKey)

	to := mailjet.RecipientsV31{}
	for _, e := range email.To {
		to = append(to, mailjet.RecipientV31{
			Email: e.Email,
			Name:  e.Name,
		})
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: email.From.Email,
				Name:  email.From.Name,
			},
			To:       &to,
			Subject:  email.Template.Subject,
			HTMLPart: email.Template.Html,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := client.SendMailV31(&messages)
	if err != nil {
		fmt.Println(fmt.Errorf("Mailjet sendmail error, %w", err))
		return err
	}
	return nil
}
