package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/s2dio-tech/mindgra-backend/common"
	_authCommon "github.com/s2dio-tech/mindgra-backend/common/auth"
	httpCommon "github.com/s2dio-tech/mindgra-backend/common/http"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
	userUsecase domain.UserUsecase
}

func InitHandlers(
	as domain.AuthUsecase,
	us domain.UserUsecase,
) *AuthHandler {
	return &AuthHandler{
		authUsecase: as,
		userUsecase: us,
	}
}

func (u *AuthHandler) Login(c *gin.Context) {
	var schema LoginRequestSchema
	// bind request context to data struct
	if err := c.Bind(&schema); err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	// validator data struct
	v := validator.New()
	if err := v.Struct(schema); err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	user, err := u.authUsecase.Authentication(c, domain.User{
		Email:    schema.Email,
		Password: schema.Password,
	})
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	AuthSuccess(u, c, *user)
}

func AuthSuccess(u *AuthHandler, c *gin.Context, user domain.User) {
	token, refreshToken, err := _authCommon.GenerateJwtTokenPair(user)
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
		},
		"token":        token,
		"refreshToken": refreshToken,
	})
}

func (u *AuthHandler) Grant(c *gin.Context) {
	var schema GrantRequestSchema
	// bind request context to data struct
	if err := c.Bind(&schema); err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	// validator data struct
	v := validator.New()
	if err := v.Struct(schema); err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	newToken, err := u.authUsecase.GrantNewAccessToken(c, schema.RefreshToken)
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

func (u *AuthHandler) ForgotPassword(c *gin.Context) {
	var schema ForgotPasswordSchema
	// bind request context to data struct
	if err := c.Bind(&schema); err != nil {
		fmt.Print(err)
		httpCommon.ErrorResponse(c, err)
		return
	}

	// validator data struct
	v := validator.New()
	if err := v.Struct(schema); err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	err := u.authUsecase.ForgotPassword(c, schema.Email)
	if err != nil {
		httpCommon.ErrorResponse(c, err)
	}

	c.JSON(http.StatusNoContent, nil)
}

func (u *AuthHandler) ResetPassword(c *gin.Context) {
	var schema ResetPasswordSchema
	// bind request context to data struct
	if err := c.Bind(&schema); err != nil {
		fmt.Print(err)
		httpCommon.ErrorResponse(c, err)
		return
	}

	// validator data struct
	v := validator.New()
	if err := v.Struct(schema); err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	// verify reset password token
	claims, ok := _authCommon.ExtractJwtClaims(schema.Token, common.AppConfig.TokenSecret)
	if !ok {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}
	var payload ResetPasswordTokenClaims
	err := common.ConvertMapToStruct(claims, &payload)
	if err != nil || payload.Action != domain.AuthActionResetPassword {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	// update password
	u.userUsecase.UpdatePassword(c, payload.UserId, schema.Password)

	c.JSON(http.StatusNoContent, nil)
}

func (u *AuthHandler) VerifyResetPasswordOTP(c *gin.Context) {
	var schema VerifyOTPSchema
	// bind request context to data struct
	if err := c.Bind(&schema); err != nil {
		fmt.Print(err)
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	// validator data struct
	v := validator.New()
	if err := v.Struct(schema); err != nil {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	user, err := u.authUsecase.VerifyOTP(c, schema.Email, schema.OTP)
	if err != nil || user == nil {
		httpCommon.ErrorResponse(c, common.ErrBadParamInput)
		return
	}

	//create jwt token for reset password
	token, err := _authCommon.GenerateJwtToken(
		map[string]interface{}{
			"action": domain.AuthActionResetPassword,
			"userId": user.Id,
		},
		common.AppConfig.TokenSecret,
		600,
	)
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
