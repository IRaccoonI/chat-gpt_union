package tokens

import (
	"encoding/json"
	"os"
	"sync"
)

type AccessToken struct {
	Tokens []string
	lock   sync.Mutex
}

func NewAccessToken(tokens []string) AccessToken {
	// Save the tokens to a file
	if _, err := os.Stat("access_tokens.json"); os.IsNotExist(err) {
		// Create the file
		file, err := os.Create("access_tokens.json")
		if err != nil {
			return AccessToken{}
		}
		defer file.Close()
	}
	file, err := os.OpenFile("access_tokens.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return AccessToken{}
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(tokens)
	if err != nil {
		return AccessToken{}
	}
	return AccessToken{
		Tokens: tokens,
	}
}

func (a *AccessToken) GetToken() string {
	a.lock.Lock()
	defer a.lock.Unlock()

	if len(a.Tokens) == 0 {
		return ""
	}

	token := a.Tokens[0]
	a.Tokens = append(a.Tokens[1:], token)
	return token
}
