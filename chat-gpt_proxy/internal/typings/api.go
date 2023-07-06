package types

type API_request_message struct {
	Messages        []API_message `json:"messages"`
	Stream          bool          `json:"stream"`
	Model           string        `json:"model"`
	ConversationID  *string       `json:"conversation_id"`
	ParentMessageID *string       `json:"parent_message_id"`
}

type API_message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type API_request_chats_item struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Create_time string `json:"create_time"`
	Update_time string `json:"update_time"`
	// mapping nil `json:"mapping"`;
	// current_node nil `json:"current_node`;
}

type API_request_chats struct {
	Items                     []API_request_chats_item `json:"items"`
	Total                     int                      `json:"total"`
	Limit                     int                      `json:"limit"`
	Offset                    int                      `json:"offset"`
	Has_missing_conversations bool                     `json:"has_missing_conversations"`
}
