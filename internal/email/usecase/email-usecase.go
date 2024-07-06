package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

type emailUsecase struct {
	sender domain.EmailSender
}

func Init(sender domain.EmailSender) domain.EmailUsecase {
	return &emailUsecase{
		sender: sender,
	}
}

func (u *emailUsecase) SendEmail(c context.Context, email domain.Email) error {
	if err := u.sender.SendEmail(email); err != nil {
		log.Fatalf("Send mail error. %v", err)
		return common.ErrInternalServerError
	}

	return nil
}

func (u *emailUsecase) GetOTPMailTemplate(name string, otp string) *domain.EmailTemplate {
	return &domain.EmailTemplate{
		Subject: "Your OTP",
		Html: fmt.Sprintf(`Hello %s<br/>
Please use the verification code below on MindGra website.<br/>
<b>%s</b><br/>
If you didn't request this, you can ignore this email.`, name, otp),
	}
}

func (u *emailUsecase) GetPasswordChangedMailTemplate(name string) *domain.EmailTemplate {
	return &domain.EmailTemplate{
		Subject: "Password Updated",
		Html: fmt.Sprintf(`Hello %s,
Your login password was updated!
If you did not make this request, please ignore this email.`, name),
	}
}
