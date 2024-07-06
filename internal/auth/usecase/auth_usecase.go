package usecase

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/s2dio-tech/mindgra-backend/common"
	_authCommon "github.com/s2dio-tech/mindgra-backend/common/auth"
	"github.com/s2dio-tech/mindgra-backend/domain"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	tokenRepo    domain.TokenRepository
	userRepo     domain.UserRepository
	emailUsecase domain.EmailUsecase
}

func InitAuthUsecase(
	tokenRepo domain.TokenRepository,
	userRepo domain.UserRepository,
	emailUsecase domain.EmailUsecase,
) domain.AuthUsecase {
	return &authUsecase{
		tokenRepo:    tokenRepo,
		userRepo:     userRepo,
		emailUsecase: emailUsecase,
	}
}

func (u *authUsecase) Authentication(c context.Context, user domain.User) (res *domain.User, err error) {
	res, err = u.userRepo.FindByEmail(user.Email)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if res == nil {
		return nil, common.ErrInvalidCredential
	}

	//verify password
	if passVerifyErr := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(user.Password)); passVerifyErr != nil {
		return nil, common.ErrInvalidCredential
	}

	return res, nil
}
func (u *authUsecase) GrantNewAccessToken(c context.Context, refreshToken string) (token *string, err error) {
	data, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return ([]byte(common.AppConfig.RefreshTokenSecret)), nil
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	claims, ok := data.Claims.(jwt.MapClaims)
	if !ok || !data.Valid {
		return nil, common.ErrBadParamInput
	}

	id := claims["id"].(string)
	role := domain.RoleMap[claims["role"].(string)]

	newToken, _ := _authCommon.GenerateJwtToken(
		map[string]interface{}{"id": id, "role": role},
		common.AppConfig.TokenSecret,
		300,
	)

	return &newToken, nil
}

func (u *authUsecase) ForgotPassword(c context.Context, mailAddress string) error {

	// find users
	user, err := u.userRepo.FindByEmail(mailAddress)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if user == nil {
		return common.ErrNotFound
	}
	// generate OTP token
	createTime := time.Now()
	token, err := totp.Generate(totp.GenerateOpts{
		Digits:      otp.DigitsSix,
		Issuer:      common.AppConfig.AppDomain,
		AccountName: mailAddress,
	})
	if err != nil {
		log.Fatal(err)
		return common.ErrInternalServerError
	}
	code, err := totp.GenerateCode(token.Secret(), createTime)
	if err != nil {
		log.Fatal(err)
		return common.ErrInternalServerError
	}

	// store token's secret
	// remove old ones
	u.tokenRepo.DeleteByTypeAndUserId(domain.TokenTypeOTP, user.Id)
	// add created one
	u.tokenRepo.Store(&domain.Token{
		Type:      domain.TokenTypeOTP,
		Token:     token.Secret(),
		UserId:    user.Id,
		CreatedAt: createTime,
	})

	// send email
	email := &domain.Email{
		From: domain.Contact{
			Name:  common.AppConfig.AppName,
			Email: common.AppConfig.AppEmail,
		},
		To: []domain.Contact{
			{Name: user.Name, Email: user.Email},
		},
		Template: *u.emailUsecase.GetOTPMailTemplate(user.Name, code),
	}

	return u.emailUsecase.SendEmail(c, *email)
}

func (u *authUsecase) VerifyOTP(c context.Context, email string, otpCode string) (*domain.User, error) {

	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		log.Fatal(err)
		return nil, common.ErrInternalServerError
	}
	if user == nil {
		return nil, common.ErrBadParamInput
	}

	// find token
	token, err := u.tokenRepo.FindOne(domain.TokenTypeOTP, user.Id)
	if err != nil {
		log.Fatal(err)
		return nil, common.ErrInternalServerError
	}
	if token == nil {
		return nil, common.ErrBadParamInput
	}

	// validate otp
	res := totp.Validate(
		otpCode,
		token.Token,
	)

	if !res {
		return nil, common.ErrBadParamInput
	}

	return user, nil
}
