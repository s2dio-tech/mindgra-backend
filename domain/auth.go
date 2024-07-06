package domain

import (
	"context"
	"time"
)

type TokenType string

const (
	TokenTypeOTP     TokenType = "otp"
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type AuthActionType string

const (
	AuthActionResetPassword AuthActionType = "resetPassword"
)

type Token struct {
	Id        string
	UserId    string
	Token     string
	Type      TokenType
	CreatedAt time.Time
}

type TokenRepository interface {
	Store(*Token) (*string, error)
	DeleteByTypeAndUserId(tType TokenType, userId string) error
	FindToken(tokenType TokenType, token string, userId string) (*Token, error)
	FindOne(tokenType TokenType, userId string) (*Token, error)
}

type AuthUsecase interface {
	Authentication(context.Context, User) (*User, error)
	GrantNewAccessToken(c context.Context, refreshToken string) (*string, error)
	ForgotPassword(context.Context, string) error
	VerifyOTP(c context.Context, email string, otp string) (*User, error)
}
