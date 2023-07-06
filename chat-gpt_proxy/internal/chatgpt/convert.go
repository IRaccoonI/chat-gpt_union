package chatgpt

import (
	typings "freechatgpt/internal/typings"
)

func ConvertAPIRequest(api_request typings.API_request_message) typings.ChatGPTRequest {
	chatgpt_request := typings.NewChatGPTRequest(api_request.Model, api_request.ConversationID, api_request.ParentMessageID)
	for _, api_message := range api_request.Messages {
		if api_message.Role == "system" {
			api_message.Role = "tool"
		}
		chatgpt_request.AddMessage(api_message.Role, api_message.Content)
	}
	return chatgpt_request
}
