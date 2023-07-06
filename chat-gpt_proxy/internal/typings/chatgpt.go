package types

import "github.com/google/uuid"

type chatgpt_message struct {
	ID      uuid.UUID       `json:"id"`
	Author  chatgpt_author  `json:"author"`
	Content chatgpt_content `json:"content"`
}

type chatgpt_content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type chatgpt_author struct {
	Role string `json:"role"`
}

type ChatGPTRequest struct {
	Action          string            `json:"action"`
	Messages        []chatgpt_message `json:"messages"`
	ParentMessageID *string           `json:"parent_message_id,omitempty"`
	Model           string            `json:"model"`
	ConversationID  *string           `json:"conversation_id"`
}

func NewChatGPTRequest(model string, conversation_id *string, parent_message_id *string) ChatGPTRequest {
	random_parent_id := uuid.NewString()

	println(conversation_id, parent_message_id)
	if parent_message_id == nil {
		return ChatGPTRequest{
			Action:          "next",
			ParentMessageID: &random_parent_id,
			Model:           model,
			ConversationID:  conversation_id,
		}
	}
	return ChatGPTRequest{
		Action:          "next",
		ParentMessageID: parent_message_id,
		Model:           model,
		ConversationID:  conversation_id,
	}
}

func (c *ChatGPTRequest) AddMessage(role string, content string) {
	c.Messages = append(c.Messages, chatgpt_message{
		ID:      uuid.New(),
		Author:  chatgpt_author{Role: role},
		Content: chatgpt_content{ContentType: "text", Parts: []string{content}},
	})
}
