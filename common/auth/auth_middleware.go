package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

// Strips 'TOKEN ' prefix from token string
func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	return strings.Replace(tok, "Bearer ", "", 1), nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(
			c.Request,
			&request.MultiExtractor{
				&request.PostExtractionFilter{
					Extractor: request.HeaderExtractor{"Authorization"},
					Filter:    stripBearerPrefixFromTokenString,
				},
				request.ArgumentExtractor{"access_token"},
			},
			func(token *jwt.Token) (interface{}, error) {
				b := ([]byte(common.AppConfig.TokenSecret))
				return b, nil
			},
		)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		c.Set("user", domain.Profile{
			Id:   claims["id"].(string),
			Role: domain.RoleMap[claims["role"].(string)],
		})
	}
}
