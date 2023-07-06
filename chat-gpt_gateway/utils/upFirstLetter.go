package utils

import "strings"

func UpFirstLetter(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}
