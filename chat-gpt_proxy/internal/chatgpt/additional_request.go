package chatgpt

import (
	"bytes"
	"encoding/json"
	"os"

	http "github.com/bogdanfinn/fhttp"
)

func GetChats(access_token string) (*http.Response, error) {
	if http_proxy != "" && len(proxies) == 0 {
		client.SetProxy(http_proxy)
	}
	// Take random proxy from proxies.txt
	if len(proxies) > 0 {
		client.SetProxy(proxies[random_int(0, len(proxies))])
	}

	apiUrl := "https://chat.openai.com/backend-api/conversations?offset=0&limit=99&order=updated"
	if API_REVERSE_PROXY != "" {
		apiUrl = API_REVERSE_PROXY
	}

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return &http.Response{}, err
	}
	// Clear cookies
	if os.Getenv("PUID") != "" {
		request.Header.Set("Cookie", "_puid="+os.Getenv("PUID")+";")
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	request.Header.Set("Accept", "*/*")
	if access_token != "" {
		request.Header.Set("Authorization", "Bearer "+access_token)
	}
	if err != nil {
		return &http.Response{}, err
	}
	response, err := client.Do(request)
	return response, err
}

func GetChat(chatId, access_token string) (*http.Response, error) {
	if http_proxy != "" && len(proxies) == 0 {
		client.SetProxy(http_proxy)
	}
	// Take random proxy from proxies.txt
	if len(proxies) > 0 {
		client.SetProxy(proxies[random_int(0, len(proxies))])
	}

	apiUrl := "https://chat.openai.com/backend-api/conversation/" + chatId
	if API_REVERSE_PROXY != "" {
		apiUrl = API_REVERSE_PROXY
	}

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return &http.Response{}, err
	}
	// Clear cookies
	if os.Getenv("PUID") != "" {
		request.Header.Set("Cookie", "_puid="+os.Getenv("PUID")+";")
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	request.Header.Set("Accept", "*/*")
	if access_token != "" {
		request.Header.Set("Authorization", "Bearer "+access_token)
	}
	if err != nil {
		return &http.Response{}, err
	}
	response, err := client.Do(request)
	return response, err
}

type TGenTitlePost struct {
	MessageID string `json:"message_id"`
}

func GenTitle(chatID, messageID, access_token string) (*http.Response, error) {
	if http_proxy != "" && len(proxies) == 0 {
		client.SetProxy(http_proxy)
	}
	// Take random proxy from proxies.txt
	if len(proxies) > 0 {
		client.SetProxy(proxies[random_int(0, len(proxies))])
	}

	apiUrl := "https://chat.openai.com/backend-api/conversation/gen_title/" + chatID
	if API_REVERSE_PROXY != "" {
		apiUrl = API_REVERSE_PROXY
	}

	data := TGenTitlePost{
		MessageID: messageID,
	}
	// Convert the struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic("json error")
	}

	request, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return &http.Response{}, err
	}
	// Clear cookies
	if os.Getenv("PUID") != "" {
		request.Header.Set("Cookie", "_puid="+os.Getenv("PUID")+";")
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	request.Header.Set("Accept", "*/*")
	if access_token != "" {
		request.Header.Set("Authorization", "Bearer "+access_token)
	}
	if err != nil {
		return &http.Response{}, err
	}
	response, err := client.Do(request)
	return response, err
}
