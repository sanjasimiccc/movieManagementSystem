package middleware

import (
	"net/http"
)

const APIKey = "my_extremely_secure_api_key_xd"

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if key := getAPIKeyFromRequest(r); key != APIKey {
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
