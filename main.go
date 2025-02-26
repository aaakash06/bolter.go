package main

import (
	"log"
	"net/http"
	"time"

	"bolter/handlers"
	"bolter/middleware"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Initialize routes
	initializeRoutes(r)

	// Start server
	log.Printf("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initializeRoutes(r *mux.Router) {
	// Create a subrouter for /api
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/template/{templateID}", handlers.TemplateHandler).Methods("GET")

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/stream")
		for i := 0; i < 10; i++ {
			w.Write([]byte("Hello, World!\n"))
			w.(http.Flusher).Flush()
			time.Sleep(1 * time.Second)
		}
		// w.Write([]byte("Hello, World!\n"))
	}).Methods("GET")

	// Add middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)
}
