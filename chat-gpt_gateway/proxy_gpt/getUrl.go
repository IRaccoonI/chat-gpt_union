package proxy_gpt

import (
	"os"
)

func GetUrl() string {
	return "http://" + os.Getenv("PROXY_HOST") + ":" + os.Getenv("PROXY_PORT")
}
