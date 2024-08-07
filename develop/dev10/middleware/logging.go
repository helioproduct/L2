package middleware

import (
	"log"
	"net/http"
)

// Обворачивает функцию next, логируя запрос
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s", req.Method, req.RequestURI)
		next.ServeHTTP(w, req)
	})
}
