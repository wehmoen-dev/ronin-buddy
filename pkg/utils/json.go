package utils

import (
	"encoding/json"
	"github.com/sethvargo/go-githubactions"
)

func Json(v interface{}) string {
	res, err := json.Marshal(v)
	if err != nil {
		githubactions.Debugf("Failed to marshal JSON: %v", err)
		return ""
	}
	return string(res)
}
