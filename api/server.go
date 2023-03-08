package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type RestApiServer struct {
	listenAddr string
}

// Creates reference to RestApiServer
func NewRestApi(listenAddr string) *RestApiServer {
	return &RestApiServer{
		listenAddr: listenAddr,
	}
}

// Runs mux router with all defined handlers
func (s *RestApiServer) Run() {
	r := mux.NewRouter()
	r.Use(corsMiddleware, loggerMiddleware)
	sub := r.PathPrefix("/api").Subrouter()

	sub.HandleFunc("/manager/{path:.*}", s.handleManager).Methods(http.MethodGet)
	sub.HandleFunc("/stream/{path:.*}", s.handleStreamFile).Methods(http.MethodGet)
	sub.HandleFunc("/download/{path:.*}", s.handleDownloadFile).Methods(http.MethodGet)

	fmt.Printf("Server running on host %s\n", s.listenAddr)

	log.Fatal(http.ListenAndServeTLS(s.listenAddr, os.Getenv("CRT_PATH"), os.Getenv("KEY_PATH"), r))
}
