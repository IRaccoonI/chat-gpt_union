package handlers

import (
	"fmt"
	"gpt-gateway/db"
	"net/http"
	"os"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type TLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TLoginResponse struct {
	Token string `json:"token"`
}

func Login(c *gin.Context) {
	var user TLoginRequest

	// Bind the JSON data from the request body to the User struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	keycloak := gocloak.NewClient(os.Getenv("KEYCLOAK_HOST"))

	// Получение токена доступа
	token, err := keycloak.Login(context.Background(), os.Getenv("KEYCLOAK_CLIENTID"), os.Getenv("KEYCLOAK_SECRET"), os.Getenv("KEYCLOAK_REALM"), user.Email, user.Password)
	if err != nil {
		fmt.Printf("Ошибка при входе: %s\n", err.Error())
		c.String(401, "Неверный логин/пароль")
		return
	}

	// Выполнение запросов к Keycloak с использованием полученного токена доступа
	// Например, получение информации о пользователе
	userInfo, err := keycloak.GetUserInfo(context.Background(), token.AccessToken, os.Getenv("KEYCLOAK_REALM"))
	if err != nil {
		fmt.Printf("Ошибка при получении информации о пользователе: %s\n", err.Error())
		c.String(401, "Jwt неверный")
		return
	}

	// Create user if not exist
	var dbUser db.User
	if result := db.DbClient.Where("username = ?", userInfo.Email).First(&dbUser); result.Error != nil {
		dbUser = db.User{
			Username: user.Email,
		}

		result := db.DbClient.Create(&dbUser)

		if result.Error != nil {
			c.String(400, "Some authorization create user error ;(")
			return
		}
	}

	var dbToken db.VerifiedToken
	if result := db.DbClient.Where("token = ?", token.AccessToken).First(&dbToken); result.Error != nil {
		dbToken = db.VerifiedToken{
			Token: token.AccessToken,
			User:  dbUser,
		}

		result := db.DbClient.Create(&dbToken)

		if result.Error != nil {
			c.String(400, "Some authorization create verified token ;(")
			return
		}
	}

	c.JSON(200, TLoginResponse{
		Token: dbToken.Token,
	})

}
