package middleware

import (
	"errors"
	"gpt-gateway/db"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// MiddlewareJWT промежуточный слой для проверки JWT
func MiddlewareJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получение JWT из куки
		cookieToken, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует token кука"})
			return
		}

		tokenString := strings.Replace(cookieToken, "Bearer ", "", 1)

		var dbToken db.VerifiedToken
		if result := db.DbClient.Where("token = ?", tokenString).Preload("User").First(&dbToken); result.Error != nil {
			// c.String(401, "Невалидный токен")
			c.AbortWithError(401, errors.New("невалидный токен"))
			return
		}

		c.Set("username", dbToken.User.Username)

		c.Next()
	}
}
