// Package httpx holds small protocol-level HTTP middlewares shared across the
// mounted APIs. It deliberately knows nothing about Astrate domains.
package httpx

import "net/http"

// CORS returns a middleware that answers preflight requests and stamps
// Access-Control headers for the allowed origins ("*" allows any). Preflights
// are answered before any downstream auth runs — browsers send them without
// an Authorization header. Requests from origins not in the list pass through
// unstamped: the browser, not the server, enforces the block.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	wildcard := false
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		if o == "*" {
			wildcard = true
			continue
		}
		allowed[o] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" || (!wildcard && !allowed[origin]) {
				next.ServeHTTP(w, r)
				return
			}

			h := w.Header()
			if wildcard {
				h.Set("Access-Control-Allow-Origin", "*")
			} else {
				h.Set("Access-Control-Allow-Origin", origin)
				h.Add("Vary", "Origin")
			}

			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				h.Set("Access-Control-Max-Age", "600")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
