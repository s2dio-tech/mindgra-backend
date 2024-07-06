package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/s2dio-tech/mindgra-backend/common/auth"
	httpCommon "github.com/s2dio-tech/mindgra-backend/common/http"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func InitHandlers(us domain.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUsecase: us,
	}
}

// @Summary
// @Schemes
// @Description user registration
// @Accept json
// @Produce json
// @Success 200
// @Router /users/ [get]
func (u *UserHandler) Register(c *gin.Context) {
	var user domain.User
	// bind request context to data struct
	if err := c.Bind(&user); err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	// validator data struct
	v := validator.New()
	if err := v.Struct(user); err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	userId, err := u.UserUsecase.Registration(c, &user)
	if err != nil {
		httpCommon.ErrorResponse(c, err)
		return
	}

	newUser := domain.User{
		Id:    *userId,
		Name:  user.Name,
		Email: user.Email,
	}
	token, refreshToken, err := auth.GenerateJwtTokenPair(newUser)
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    newUser.Id,
			"name":  newUser.Name,
			"email": newUser.Email,
		},
		"token":        token,
		"refreshToken": refreshToken,
	})
}
