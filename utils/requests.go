package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CloneResponse struct {
	IsLua bool `json:"isLua"`
}

func CheckIsLua(repo string) (bool, error) {
	url := "http://localhost:1512/clone"

	data := map[string]string{"Repo": repo}
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
