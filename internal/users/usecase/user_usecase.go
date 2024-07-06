package usecase

import (
	"context"
	"time"

	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/domain"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo     domain.UserRepository
	emailUsecase domain.EmailUsecase
}

func InitUserUsecase(
	repo domain.UserRepository,
	emailUsecase domain.EmailUsecase,
) domain.UserUsecase {
	return &userUsecase{
		userRepo:     repo,
		emailUsecase: emailUsecase,
	}
}

func hashPassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", common.ErrInternalServerError
	}
	return string(hash), nil
}

func (u *userUsecase) Registration(c context.Context, user *domain.User) (res *string, err error) {
	// check email exists
	oldUser, err := u.userRepo.FindByEmail(user.Email)
	if err != nil {
		return
	}
	if oldUser != nil {
		return nil, common.ErrEmailDuplicate
	}

	hashPassword, err := hashPassword(user.Password)
	if err != nil {
		return
	}
	// insert to db
	userId, err := u.userRepo.Create(&domain.User{
		Name:      user.Name,
		Email:     user.Email,
		Password:  hashPassword,
		Role:      domain.RoleMember,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return nil, common.ErrInternalServerError
	}

	return userId, nil
}

func (u *userUsecase) FindByEmail(c context.Context, email string) (user *domain.User, err error) {
	user, err = u.userRepo.FindByEmail(email)
	if err != nil {
		return
	}
	if user == nil {
		return nil, common.ErrNotFound
	}

	return
}

func (u *userUsecase) FindById(c context.Context, id string) (user *domain.User, err error) {
	user, err = u.userRepo.FindById(id)
	if err != nil {
		return
	}
	if user == nil {
		return nil, common.ErrNotFound
	}

	return
}

func (u *userUsecase) UpdatePassword(c context.Context, userId string, password string) error {

	// find user
	user, err := u.userRepo.FindById(userId)
	if err != nil {
		return common.ErrInternalServerError
	}
	if user == nil {
		return common.ErrNotFound
	}

	// update password
	hashPassword, err := hashPassword(password)
	if err != nil {
		return common.ErrInternalServerError
	}
	err = u.userRepo.Update(userId, map[string]interface{}{
		"password": hashPassword,
	})

	if err != nil {
		return common.ErrInternalServerError
	}

	// send email
	mail := &domain.Email{
		From: domain.Contact{
			Name:  common.AppConfig.AppName,
			Email: common.AppConfig.AppEmail,
		},
		To: []domain.Contact{
			{
				Email: user.Email,
				Name:  user.Name,
			},
		},
		Template: *u.emailUsecase.GetPasswordChangedMailTemplate(user.Name),
	}

	return u.emailUsecase.SendEmail(c, *mail)
}
