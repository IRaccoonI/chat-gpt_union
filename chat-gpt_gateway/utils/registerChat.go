package utils

import (
	db "gpt-gateway/db"
	"gpt-gateway/proxy_gpt"
)

func RegisterConversation(chatID string, username string) string {
	var user db.User
	db.DbClient.First(&user, "username = ?", username)

	var chat db.Chat
	db.DbClient.First(&chat, "id = ?", chatID)
	// Conversation already exist
	if chat.ID != "" {
		return ""
	}

	originalChat := proxy_gpt.GetChat(chatID)

	newChat := db.Chat{
		ID:   chatID,
		Name: originalChat.Title,
		User: user,
	}

	db.DbClient.Create(&newChat)

	return originalChat.Title
}
