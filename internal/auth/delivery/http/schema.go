package http

import "github.com/s2dio-tech/mindgra-backend/domain"

type LoginRequestSchema struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type GrantRequestSchema struct {
	RefreshToken string `json:"refreshToken" validate:"required,max=512"`
}

type ForgotPasswordSchema struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type VerifyOTPSchema struct {
	OTP   string `json:"otp" validate:"required,min=6,max=6"`
	Email string `json:"email" validate:"required,email,max=255"`
}

type ResetPasswordSchema struct {
	Token    string `json:"token" validate:"required"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type ResetPasswordTokenClaims struct {
	Action domain.AuthActionType `json:"action"`
	UserId string                `json:"userId"`
	exp    int64                 `json:"exp"`
}
