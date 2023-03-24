package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arekbor/file-manager-server/utils"
)

// Allows cors origin
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware", r.Method)

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers:", "Origin, Content-Type, X-Auth-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
		log.Println("Executing middleware again")
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\033[33mREQUEST:\033[0m\033[32m [Proto: %s], [Host: %s], [Method: %s], [Path: %s]\033[0m\n", utils.GetProtoFromRequest(r), r.Host, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
