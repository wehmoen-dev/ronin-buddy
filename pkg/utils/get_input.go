package utils

import "github.com/sethvargo/go-githubactions"

func GetActionInputOrDefault(key string, defaultValue string) string {
	value := githubactions.GetInput(key)
	if value == "" {
		return defaultValue
	}
	return value
}
