package handlers

import (
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func TestAuth(c *gin.Context) {
	keycloak := gocloak.NewClient("http://localhost:8080")

	// Получение токена доступа
	token, err := keycloak.Login(context.Background(), "smth", "Y0eGDgSlcjGs1i74DmhqIGINF1pdpEQc", "master", "test", "test")
	if err != nil {
		fmt.Printf("Ошибка при входе: %s\n", err.Error())
		c.String(500, "ups")
		return
	}

	// Выполнение запросов к Keycloak с использованием полученного токена доступа
	// Например, получение информации о пользователе
	userInfo, err := keycloak.GetUserInfo(context.Background(), token.AccessToken, "master")
	if err != nil {
		fmt.Printf("Ошибка при получении информации о пользователе: %s\n", err.Error())
		c.String(500, "ups 2")
		return
	}

	fmt.Printf("Имя пользователя: %s\n", userInfo.PreferredUsername)
	c.JSON(200, token.AccessToken)
}
