package middleware

import (
	"net/http"

	"github.com/gauravst/real-time-chat/internal/config"
)

func CORS(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow only the specific frontend URL
			w.Header().Set("Access-Control-Allow-Origin", cfg.ClientUrl)

			// Allow credentials (IMPORTANT for cookies/tokens)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Allowed methods
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			// Allowed headers
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
