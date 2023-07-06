package proxy_gpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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

func SendMessage(text, conversationID, parentMessageID string) ChatCompletion {
	url := GetUrl() + "/v1/chat/completions"
	payload := APIRequest{
		Model: "text-davinci-002-render-sha",
		Messages: []api_message{
			{
				Role:    "user",
				Content: text,
			},
		},
	}

	if conversationID != "" {
		payload.ConversationId = &conversationID
	}

	if parentMessageID != "" {
		payload.ParentMessageId = &parentMessageID
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		panic("Error marshaling JSON:")
	}

	// Create a new HTTP GET request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payloadJson)))
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

	var bodyJson ChatCompletion
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		panic("Proxy error")
	}

	return bodyJson
}

type ChatGPTResponse struct {
	Message        Message     `json:"message"`
	ConversationID string      `json:"conversation_id"`
	Error          interface{} `json:"error"`
}

func SendMessageStream(c *gin.Context, text, conversationID, parentMessageID string) <-chan ChatGPTResponse {
	ch := make(chan ChatGPTResponse)

	url := GetUrl() + "/v1/chat/completions"
	payload := APIRequest{
		Stream: true,
		Model:  "text-davinci-002-render-sha",
		Messages: []api_message{
			{
				Role:    "user",
				Content: text,
			},
		},
	}

	if conversationID != "" {
		payload.ConversationId = &conversationID
	}

	if parentMessageID != "" {
		payload.ParentMessageId = &parentMessageID
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		panic("Error marshaling JSON:")
	}

	// Create a new HTTP GET request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payloadJson)))
	if err != nil {
		panic("Error creating request:")
	}

	// Send the request using the default HTTP client

	go func() {
		defer close(ch)

		client := http.DefaultClient
		response, err := client.Do(req)
		if err != nil {
			panic("Error sending request:")
		}
		defer response.Body.Close()

		reader := bufio.NewReader(response.Body)

		for {
			originalLine, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return
			}

			if len(originalLine) < 6 {
				continue
			}

			line := originalLine[6:]

			if strings.HasPrefix(line, "[DONE]") {
				break
			}

			var original_response ChatGPTResponse
			err = json.Unmarshal([]byte(line), &original_response)
			if err != nil {
				continue
			}

			ch <- original_response
		}
	}()

	return ch
}

type Choices struct {
	Delta        Delta       `json:"delta"`
	Index        int         `json:"index"`
	FinishReason interface{} `json:"finish_reason"`
}

type Delta struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}

type ChatCompletionChunk struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
}

func (chunk *ChatCompletionChunk) String() string {
	resp, _ := json.Marshal(chunk)
	return string(resp)
}

func StopChunk(reason string) ChatCompletionChunk {
	return ChatCompletionChunk{
		ID:      "chatcmpl-QXlha2FBbmROaXhpZUFyZUF3ZXNvbWUK",
		Object:  "chat.completion.chunk",
		Created: 0,
		Model:   "gpt-3.5-turbo-0301",
		Choices: []Choices{
			{
				Index:        0,
				FinishReason: reason,
			},
		},
	}
}
