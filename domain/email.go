package domain

import (
	"context"
)

type Contact struct {
	Email string
	Name  string
}

type EmailTemplate struct {
	Subject string
	Html    string
}

type Email struct {
	From     Contact
	To       []Contact
	Cc       *[]Contact
	Template EmailTemplate
}

type EmailSender interface {
	SendEmail(Email) error
}

type EmailUsecase interface {
	SendEmail(context.Context, Email) error
	GetOTPMailTemplate(string, string) *EmailTemplate
	GetPasswordChangedMailTemplate(string) *EmailTemplate
}
