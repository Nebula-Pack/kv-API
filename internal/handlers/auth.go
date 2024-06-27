package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Nebula-Pack/kv-API/internal/auth"
)

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AdminLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// For simplicity, we are hardcoding the first admin credentials
		if req.Username != "KeaganGilmore" || req.Password != "#Keagan4200" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(req.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
