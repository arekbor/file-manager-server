package api

import (
	"fmt"
	"log"
	"net/http"

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

	sub.HandleFunc("/upload", s.handleUpload).Methods(http.MethodPost, http.MethodOptions)
	sub.HandleFunc("/manager/{path:.*}", s.handleManager).Methods(http.MethodGet)
	sub.HandleFunc("/stream/{path:.*}", s.handleStreamFile).Methods(http.MethodGet)
	sub.HandleFunc("/download/{path:.*}", s.handleDownloadFile).Methods(http.MethodGet)

	fmt.Printf("Server running on host %s\n", s.listenAddr)

	server := &http.Server{
		Addr:    s.listenAddr,
		Handler: r,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
		return
	}
}
