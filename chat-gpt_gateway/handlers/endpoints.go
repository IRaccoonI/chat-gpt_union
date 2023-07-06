package handlers

import "github.com/gin-gonic/gin"

type TEndpoints struct {
	AzureOpenAI    bool            `json:"azureOpenAI"`
	BingAI         bool            `json:"bingAI"`
	ChatGptBrowser IEndpointsModel `json:"chatGptBrowser"`
	OpenAI         IEndpointsModel `json:"openAI"`
}

type IEndpointsModel struct {
	AvailableModels []string `json:"availableModels"`
}

func GetEndpoints(c *gin.Context) {
	// username := c.GetString("username")

	c.JSON(200, TEndpoints{
		AzureOpenAI:    false,
		BingAI:         false,
		ChatGptBrowser: IEndpointsModel{AvailableModels: []string{"text-davinci-002-render-sha"}},
		OpenAI:         IEndpointsModel{AvailableModels: []string{}},
	})
}
