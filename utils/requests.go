package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type CloneResponse struct {
	GithubURL    string            `json:"github_url"`
	IsLua        bool              `json:"isLua"`
	HasRockspec  bool              `json:"hasRockspec,omitempty"`
	ScanResponse map[string]string `json:"scanResponse,omitempty"`
	Version      string            `json:"version"` // Include the version field
}

func CheckIsLua(repo, version string) (CloneResponse, error) {
	url := "http://localhost:1912/clone"

	data := map[string]string{"repo": repo}
	if version != "" {
		data["version"] = version
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return CloneResponse{}, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return CloneResponse{}, err
	}
	defer resp.Body.Close()

	var cloneResp CloneResponse
	if err := json.NewDecoder(resp.Body).Decode(&cloneResp); err != nil {
		return CloneResponse{}, err
	}

	return cloneResp, nil
}
