package handlers

import (
	"gpt-gateway/db"
	"gpt-gateway/proxy_gpt"
	utils "gpt-gateway/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TGetConversationsResponse struct {
	Conversations []TConversation `json:"conversations"`
	PageNumber    int             `json:"pageNumber"`
	PageSize      int             `json:"pageSize"`
	Pages         int             `json:"pages"`
}

type TConversation struct {
	ConversationID          string         `json:"conversationId"`
	Title                   string         `json:"title"`
	User                    string         `json:"user,omitempty"`
	Endpoint                EModelEndpoint `json:"endpoint"`
	Suggestions             []string       `json:"suggestions,omitempty"`
	Messages                []TMessage     `json:"messages,omitempty"`
	Tools                   []interface{}  `json:"tools,omitempty"`
	CreatedAt               string         `json:"createdAt"`
	UpdatedAt               string         `json:"updatedAt"`
	ModelLabel              string         `json:"modelLabel,omitempty"`
	Examples                []TExample     `json:"examples,omitempty"`
	ChatGptLabel            string         `json:"chatGptLabel,omitempty"`
	UserLabel               string         `json:"userLabel,omitempty"`
	Model                   string         `json:"model,omitempty"`
	PromptPrefix            string         `json:"promptPrefix,omitempty"`
	Temperature             float64        `json:"temperature,omitempty"`
	TopP                    float64        `json:"topP,omitempty"`
	TopK                    int            `json:"topK,omitempty"`
	Context                 string         `json:"context,omitempty"`
	TopPAlternative         float64        `json:"top_p,omitempty"`
	PresencePenalty         float64        `json:"presence_penalty,omitempty"`
	Jailbreak               bool           `json:"jailbreak,omitempty"`
	JailbreakConversationId string         `json:"jailbreakConversationId,omitempty"`
	ConversationSignature   string         `json:"conversationSignature,omitempty"`
	ParentMessageId         string         `json:"parentMessageId,omitempty"`
	ClientId                string         `json:"clientId,omitempty"`
	InvocationId            string         `json:"invocationId,omitempty"`
	FrequencyPenalty        float64        `json:"frequency_penalty,omitempty"`
	ToneStyle               string         `json:"toneStyle,omitempty"`
}

type TExample struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type EModelEndpoint string

const (
	AzureOpenAI    EModelEndpoint = "azureOpenAI"
	OpenAI         EModelEndpoint = "openAI"
	BingAI         EModelEndpoint = "bingAI"
	ChatGPT        EModelEndpoint = "chatGPT"
	ChatGPTBrowser EModelEndpoint = "chatGPTBrowser"
	Google         EModelEndpoint = "google"
	GptPlugins     EModelEndpoint = "gptPlugins"
)

func GetConversations(c *gin.Context) {
	username := c.GetString("username")

	var user db.User
	db.DbClient.First(&user, "username = ?", username)

	var chats []db.Chat
	db.DbClient.Find(&chats, "user_id = ?", user.ID)

	convos := []TConversation{}

	for _, chat := range chats {
		convos = append(convos, TConversation{
			ConversationID:   chat.ID,
			Title:            chat.Name,
			Endpoint:         "openAI",
			FrequencyPenalty: 0,
			Model:            "gpt-3.5-turbo",
			PresencePenalty:  0,
			Temperature:      1,
			TopP:             1,
		})
	}

	c.JSON(200, TGetConversationsResponse{
		Conversations: convos,
		PageNumber:    0,
		PageSize:      20,
		Pages:         0,
	})
}

func GetConversation(c *gin.Context) {
	username := c.GetString("username")
	var user db.User
	db.DbClient.First(&user, "username = ?", username)

	chatID, success := c.Params.Get("id")
	if !success {
		c.String(404, "Id wrong")
		return
	}

	var dbChat db.Chat
	db.DbClient.First(&dbChat, "user_id = ? AND id = ?", user.ID, chatID)
	if dbChat.ID == "" {
		c.String(403, "Permission denied")
		return
	}

	messages := getMessages(chatID)

	chat := TConversation{
		ConversationID:   chatID,
		Title:            "test Title",
		Endpoint:         "openAI",
		FrequencyPenalty: 0,
		Model:            "gpt-3.5-turbo",
		PresencePenalty:  0,
		Temperature:      1,
		TopP:             1,
		Messages:         messages,
	}

	var chats []db.Chat
	db.DbClient.Find(&chats, "user_id = ?", user.ID)

	c.JSON(200, chat)
}

func GetMessages(c *gin.Context) {
	username := c.GetString("username")
	var user db.User
	db.DbClient.First(&user, "username = ?", username)

	chatID, success := c.Params.Get("id")
	if !success {
		c.String(404, "Id wrong")
		return
	}

	var dbChat db.Chat
	db.DbClient.First(&dbChat, "user_id = ? AND id = ?", user.ID, chatID)
	if dbChat.ID == "" {
		c.String(403, "Permission denied")
		return
	}

	c.JSON(200, getMessages(chatID))
}

type TDeleteConversationRequest struct {
	Arg struct {
		ConversationId string `json:"conversationId,omitempty"`
		Source         string `json:"source,omitempty"`
	} `json:"arg"`
}

func DeleteConversation(c *gin.Context) {
	var payload TDeleteConversationRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Arg.ConversationId != "" {
		var chat db.Chat
		db.DbClient.Where("id = ?", payload.Arg.ConversationId).First(&chat)
		db.DbClient.Delete(&chat)
	} else {
		var chats []db.Chat
		db.DbClient.Find(&chats)
		for _, chat := range chats {
			db.DbClient.Delete(&chat)
		}
	}

	c.String(200, "nu OK")
}

func getMessages(chatID string) []TMessage {
	var messages []TMessage
	originalChat := proxy_gpt.GetChat(chatID)

	for messageID, message := range originalChat.Mapping {
		var messageText string

		if message.Message.Content.Parts != nil && len(message.Message.Content.Parts) > 0 {
			messageText = message.Message.Content.Parts[0]
		}

		messages = append(messages, TMessage{
			ConversationID:  chatID,
			MessageID:       messageID,
			ParentMessageID: message.Parent,
			Text:            messageText,
			Sender:          utils.UpFirstLetter(message.Message.Author.Role),
			CreateTime:      message.Message.CreateTime,
		})
	}

	return messages
}
