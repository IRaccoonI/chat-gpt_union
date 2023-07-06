package proxy_gpt

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Chats struct {
	Items                   []ChatItem
	Total                   int
	Limit                   int
	Offset                  int
	HasMissingConversations bool
}

type ChatItem struct {
	ID          string
	Title       string
	CreateTime  string
	UpdateTime  string
	Mapping     interface{}
	CurrentNode interface{}
}

func GetChats(chatID string) Chats {
	url := GetUrl() + "/v1/chat/conversations"

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("Error creating request:")
	}

	// Send the request using the default HTTP client
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic("Error sending request:")
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Error reading response body:")
	}

	var bodyJson Chats
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		panic("Proxy error")
	}

	return bodyJson
}
