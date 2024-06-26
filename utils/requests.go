// utils/request.go
package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CloneResponse struct {
	IsLua bool `json:"isLua"`
}

func CheckIsLua(key, value string) (bool, error) {
	url := "http://localhost:4242/clone"

	data := map[string]string{"key": key, "value": value}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var cloneResp CloneResponse
	if err := json.NewDecoder(resp.Body).Decode(&cloneResp); err != nil {
		return false, err
	}

	return cloneResp.IsLua, nil
}
