package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/s2dio-tech/mindgra-backend/common"
)

func ErrorResponse(c *gin.Context, err error) {
	switch err {
	case common.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		break
	case common.ErrUnauthentication:
	case common.ErrInvalidCredential:
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		break
	case common.ErrUnauthorization:
		c.JSON(http.StatusForbidden, gin.H{
			"message": err.Error(),
		})
		break
	case common.ErrInternalServerError:
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		break
	case common.ErrBadParamInput:
	case common.ErrEmailDuplicate:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		break
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": common.ErrInternalServerError.Error(),
		})
		break
	}
}
