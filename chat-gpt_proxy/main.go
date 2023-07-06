package main

import (
	"freechatgpt/internal/tokens"
	"log"
	"os"

	"github.com/acheong08/endless"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

var HOST string
var PORT string
var ACCESS_TOKENS tokens.AccessToken

func init() {
	godotenv.Load()
	if os.Getenv("PROXY_GPT_EMAIL") == "" {
		log.Fatal("Отсутствуют переменные среды")
		return
	}

	HOST = os.Getenv("PROXY_HOST")
	PORT = os.Getenv("PROXY_PORT")
	if HOST == "" {
		HOST = "0.0.0.0"
	}
	if PORT == "" {
		PORT = "3000"
	}

	// accessToken := os.Getenv("ACCESS_TOKENS")
	// if accessToken != "" {
	// 	accessTokens := strings.Split(accessToken, ",")
	// 	ACCESS_TOKENS = tokens.NewAccessToken(accessTokens)
	// }
	// // Check if access_tokens.json exists
	// if _, err := os.Stat("access_tokens.json"); os.IsNotExist(err) {
	// 	// Create the file
	// 	file, err := os.Create("access_tokens.json")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer file.Close()
	// } else {
	// Load the tokens
	// file, err := os.Open("access_tokens.json")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// decoder := json.NewDecoder(file)
	// var token_list []string
	// err = decoder.Decode(&token_list)
	// if err != nil {
	// 	return
	// }
	token := GenToken()
	if token == "" {
		log.Fatal("Неверный логин/пароль от GPT или нет впн, одно из двух")
		return
	}
	ACCESS_TOKENS = tokens.AccessToken{
		Tokens: []string{GenToken()},
	}
	// }
}

func main() {
	router := gin.Default()

	router.Use(cors)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// admin_routes := router.Group("/admin")
	// admin_routes.Use(adminCheck)

	/// Admin routes
	// admin_routes.PATCH("/password", passwordHandler)
	// admin_routes.PATCH("/tokens", tokensHandler)
	// admin_routes.PATCH("/puid", puidHandler)
	// admin_routes.PATCH("/openai", openaiHandler)
	/// Public routes
	router.OPTIONS("/v1/chat/completions", optionsHandler)
	router.POST("/v1/chat/completions", Authorization, nightmare)
	// router.GET("/v1/engines", Authorization, engines_handler)
	// router.GET("/v1/models", Authorization, engines_handler)

	router.GET("/v1/chat/conversation/:id", Authorization, getChat)
	router.GET("/v1/chat/conversations", Authorization, getChats)

	endless.ListenAndServe(HOST+":"+PORT, router)
}
