package middleware

import (
	"net/http"

	"github.com/sanjasimiccc/movieManagementSystem/pkg/config"
)

//const APIKey = "my_extremely_secure_api_key_xd"

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if key := getAPIKeyFromRequest(r); key != config.Envs.AuthAPIKey {
			http.Error(w, "Unauthorized: invalid API key", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getAPIKeyFromRequest(r *http.Request) string {
	apiKey := r.Header.Get("X-API-Key")

	if apiKey != "" {
		return apiKey
	}
	return ""
}
