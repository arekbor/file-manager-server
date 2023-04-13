package api

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/cors"
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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("CORS_URL")},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type", "Depth", "User-Agent", "X-File-Size", "X-Requested-With", "If-Modified-Since", "X-File-Name", "Cache-Control"},
		AllowCredentials: true,
	}), loggerMiddleware)

	sub := r.PathPrefix("/api").Subrouter()

	sub.HandleFunc("/upload", s.handleUpload).Methods(http.MethodPost, http.MethodOptions)
	sub.HandleFunc("/manager/{path:.*}", s.handleManager).Methods(http.MethodGet)
	sub.HandleFunc("/stream/{path:.*}", s.handleStreamFile).Methods(http.MethodGet)
	sub.HandleFunc("/download/{path:.*}", s.handleDownloadFile).Methods(http.MethodGet)
	sub.HandleFunc("/folderNames", s.handleGetAllFolderNames).Methods(http.MethodGet)

	log.Printf("Server running on host %s\n", s.listenAddr)

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
