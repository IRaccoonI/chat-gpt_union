// https://chat.openai.com/backend-api/conversation/

package main

import (
	"encoding/json"
	chatgpt "freechatgpt/internal/chatgpt"

	"github.com/gin-gonic/gin"
)

type Chat struct {
	Title             string        `json:"title"`
	CreateTime        float64       `json:"create_time"`
	UpdateTime        float64       `json:"update_time"`
	Mapping           Mapping       `json:"mapping"`
	ModerationResults []interface{} `json:"moderation_results"`
	CurrentNode       string        `json:"current_node"`
}

type Mapping map[string]Item

type FinishDetails struct {
	Type string `json:"type"`
	Stop string `json:"stop"`
}

type MessageMetadata struct {
	ModelSlug     string        `json:"model_slug"`
	FinishDetails FinishDetails `json:"finish_details"`
	Timestamp     string        `json:"timestamp_"`
}

type Item struct {
	ID       string   `json:"id"`
	Message  Message  `json:"message"`
	Parent   string   `json:"parent"`
	Children []string `json:"children"`
}

type Message struct {
	ID         string          `json:"id"`
	Author     Author          `json:"author"`
	CreateTime float64         `json:"create_time"`
	Content    Content         `json:"content"`
	Status     string          `json:"status"`
	Weight     float64         `json:"weight"`
	Metadata   MessageMetadata `json:"metadata"`
	EndTurn    bool            `json:"end_turn"`
	Recipient  string          `json:"recipient"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type Author struct {
	Role     string   `json:"role"`
	Metadata struct{} `json:"metadata"`
}

func getChat(c *gin.Context) {
	var chatId string
	chatId, success := c.Params.Get("id")
	if !success {
		c.String(404, "Id wrong")
		return
	}

	token := ACCESS_TOKENS.GetToken()

	response, err := chatgpt.GetChat(chatId, token)

	if err != nil {
		c.String(404, "Chat not found")
		return
	}

	defer response.Body.Close()
	if chatgpt.Handle_request_error(c, response) {
		return
	}

	decoder := json.NewDecoder(response.Body)
	var resJson Chat
	err = decoder.Decode(&resJson)
	if err != nil {
		panic(err)
	}

	c.JSON(200, resJson)
}

type Chats struct {
	Items                   []ChatsItem `json:"items"`
	Total                   int         `json:"total"`
	Limit                   int         `json:"limit"`
	Offset                  int         `json:"offset"`
	HasMissingConversations bool        `json:"has_missing_conversations"`
}

type ChatsItem struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

func getChats(c *gin.Context) {
	token := ACCESS_TOKENS.GetToken()

	response, err := chatgpt.GetChats(token)

	if err != nil {
		c.String(404, "Chat not found")
		return
	}

	defer response.Body.Close()
	if chatgpt.Handle_request_error(c, response) {
		return
	}

	decoder := json.NewDecoder(response.Body)
	var resJson Chats
	err = decoder.Decode(&resJson)
	if err != nil {
		panic(err)
	}

	c.JSON(200, resJson)
}
