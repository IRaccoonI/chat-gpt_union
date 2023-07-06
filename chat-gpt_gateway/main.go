package main

import (
	db "gpt-gateway/db"
	handlers "gpt-gateway/handlers"
	middleware "gpt-gateway/handlers/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dbCloser := db.ConnectDb()
	defer dbCloser()
	db.AutoMigrateDb()

	// Create a new Gin router
	router := gin.Default()

	// Register the handler function
	router.POST("/api/auth/login", handlers.Login)

	router.GET("/ping", middleware.MiddlewareJWT(), handlers.Ping)
	router.GET("/api/convos", middleware.MiddlewareJWT(), handlers.GetConversations)
	router.GET("/api/convos/:id", middleware.MiddlewareJWT(), handlers.GetConversation)
	router.GET("/api/messages/:id", middleware.MiddlewareJWT(), handlers.GetMessages)
	router.GET("/api/endpoints", middleware.MiddlewareJWT(), handlers.GetEndpoints)
	router.POST("/api/convos/clear", middleware.MiddlewareJWT(), handlers.DeleteConversation)

	router.POST("/api/ask/openAI", middleware.MiddlewareJWT(), handlers.AskOpenAI)

	router.POST("/api/testAuth", handlers.TestAuth)

	HOST := os.Getenv("GATEWAY_HOST")
	PORT := os.Getenv("GATEWAY_PORT")
	if HOST == "" {
		HOST = "0.0.0.0"
	}
	if PORT == "" {
		PORT = "3000"
	}

	// Start the server
	router.Run(HOST + ":" + PORT)
}
