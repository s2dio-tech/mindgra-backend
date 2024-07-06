package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/s2dio-tech/mindgra-backend/common"
	"github.com/s2dio-tech/mindgra-backend/domain"
)

func ExtractUser(c *gin.Context) domain.Profile {
	u, _ := c.Get("user")
	if u == nil {
		panic(common.ErrInternalServerError)
	}
	return u.(domain.Profile)
}

// A Util function to generate jwt_token which can be used in the request header
// exp: seconds
func GenerateJwtToken(payload map[string]interface{}, secret string, exp int) (string, error) {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Second * time.Duration(exp)).Unix(),
	}
	for k, v := range payload {
		claims[k] = v
	}
	jwt_token.Claims = claims
	// Sign and get the complete encoded token as a string
	return jwt_token.SignedString([]byte(secret))
}

func GenerateJwtTokenPair(user domain.User) (token *string, refreshToken *string, err error) {
	payload := map[string]interface{}{
		"id":   user.Id,
		"role": user.Role,
	}
	tokenStr, _ := GenerateJwtToken(payload, common.AppConfig.TokenSecret, 60*60)                     // 1hour
	refreshTokenStr, _ := GenerateJwtToken(payload, common.AppConfig.RefreshTokenSecret, 90*24*60*60) // 90 days
	return &tokenStr, &refreshTokenStr, nil
}

func ExtractJwtClaims(tokenStr string, secret string) (map[string]interface{}, bool) {
	hmacSecret := []byte(secret)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		return nil, false
	}
}
