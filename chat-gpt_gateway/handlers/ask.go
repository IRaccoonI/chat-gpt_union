package handlers

import (
	"encoding/json"
	proxy_gpt "gpt-gateway/proxy_gpt"
	utils "gpt-gateway/utils"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type TMessage struct {
	ConversationID  string  `json:"conversationId"`
	MessageID       string  `json:"messageId"`
	NewMessageID    string  `json:"newMessageId"`
	ParentMessageID string  `json:"parentMessageId"`
	Text            string  `json:"text"`
	Sender          string  `json:"sender"`
	Unfinished      bool    `json:"unfinished"`
	Cancelled       bool    `json:"cancelled"`
	Error           bool    `json:"error"`
	IsCreatedByUser bool    `json:"isCreatedByUser"`
	CreateTime      float64 `json:"createTime"`
}

type TResponseConversation struct {
	ConversationID string `json:"conversationId"`
	Title          string `json:"title"`
}

type TResponse struct {
	Conversation    TResponseConversation `json:"conversation,omitempty"`
	ResponseMessage TMessage              `json:"responseMessage,omitempty"`
	RequestMessage  TMessage              `json:"requestMessage,omitempty"`
	Text            string                `json:"text,omitempty"`
	Final           bool                  `json:"final"`
}

type APIRequest struct {
	Messages        []api_message `json:"messages"`
	Stream          bool          `json:"stream"`
	Model           string        `json:"model"`
	PluginIDs       []string      `json:"plugin_ids"`
	ParentMessageId *string       `json:"parent_message_id"`
	ConversationId  *string       `json:"conversation_id"`
}

type api_message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletion struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Usage   usage  `json:"usage"`
	Message Msg    `json:"message"`
}
type Msg struct {
	Role            string `json:"role"`
	Content         string `json:"content"`
	ConversationID  string `json:"conversation_id"`
	MessageID       string `json:"id"`
	ParentMessageID string `json:"parent_message_id"`
}
type Choice struct {
	Index        int         `json:"index"`
	Message      Msg         `json:"message"`
	FinishReason interface{} `json:"finish_reason"`
}
type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

var mutex sync.Mutex

func AskOpenAI(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // Disable buffering for nginx

	username := c.GetString("username")
	var requestMessage TMessage

	// Bind the JSON data from the request body to the User struct
	if err := c.ShouldBindJSON(&requestMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	events := proxy_gpt.SendMessageStream(c, requestMessage.Text, requestMessage.ConversationID, requestMessage.ParentMessageID)

	var lastEvent proxy_gpt.ChatGPTResponse
	for event := range events {
		if event.Message.Content.Parts == nil || len(event.Message.Content.Parts) == 0 {
			continue
		}
		lastEvent = event

		response := TResponse{
			Text: event.Message.Content.Parts[0],
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			panic("uuuuu")
		}

		c.Writer.WriteString("data: " + string(jsonResponse) + "\n\n")
		c.Writer.Flush()
	}

	title := utils.RegisterConversation(lastEvent.ConversationID, username)

	finalResponse := TResponse{
		Final:          true,
		RequestMessage: requestMessage,
		ResponseMessage: TMessage{
			MessageID:       lastEvent.Message.ID,
			Text:            lastEvent.Message.Content.Parts[0],
			Sender:          lastEvent.Message.Author.Role,
			ParentMessageID: requestMessage.MessageID,
			ConversationID:  lastEvent.ConversationID,
		},
		Conversation: TResponseConversation{
			ConversationID: lastEvent.ConversationID,
			Title:          title,
		},
	}
	finalResponse.RequestMessage.ConversationID = lastEvent.ConversationID

	jsonFinalResponse, err := json.Marshal(finalResponse)
	if err != nil {
		panic("uuuuu")
	}

	c.Writer.WriteString("data: " + string(jsonFinalResponse) + "\n\n")
	c.String(200, "")
}
