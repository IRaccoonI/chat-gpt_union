package proxy_gpt

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Chat struct {
	Title             string
	CreateTime        int64
	UpdateTime        int64
	Mapping           Mapping
	ModerationResults []interface{}
	CurrentNode       string
}

type Mapping map[string]Item

type FinishDetails struct {
	Type string `json:"type"`
	Stop string `json:"stop"`
}

type MessageMetadata struct {
	Timestamp     string         `json:"timestamp_"`
	MessageType   string         `json:"message_type"`
	FinishDetails *FinishDetails `json:"finish_details"`
	ModelSlug     string         `json:"model_slug"`
	Recipient     string         `json:"recipient"`
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

func GetChat(chatID string) Chat {
	url := GetUrl() + "/v1/chat/conversation/" + chatID

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

	var bodyJson Chat
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		panic("Proxy error")
	}

	return bodyJson
}
