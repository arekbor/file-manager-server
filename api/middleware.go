package api

import (
	"fmt"
	"net/http"

	"github.com/arekbor/file-manager-server/utils"
)

// Allows cors origin
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\033[33mREQUEST:\033[0m\033[32m [Proto: %s], [Host: %s], [Method: %s], [Path: %s]\033[0m\n", utils.GetProtoFromRequest(r), r.Host, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
