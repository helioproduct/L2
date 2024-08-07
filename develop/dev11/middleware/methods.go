package middleware

import (
	"net/http"
)

// Обворачивает функцию next, вызывая ее только, если req.Method == method
func WithMethod(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == method {
			next.ServeHTTP(w, req)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{error:"method not allowed"}`))
		}
	})
}
