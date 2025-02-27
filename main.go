package main

import (
	"bolter/handlers"
	"bolter/middleware"
	"strconv"

	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	println("starting...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// utils.Chat()
	// getting env variables
	// key := os.Getenv("OPENAI_API_KEY")
	// fmt.Printf("godotenv : %s = %s \n", "KEY", key)
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

	api.HandleFunc("/template", handlers.TemplateHandler).Methods("GET")
	api.HandleFunc("/chat", handlers.Chat).Methods("GET")

	api.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		params := mux.Vars(r)
		w.Header().Set("Content-Type", "text/stream")
		num, err := strconv.Atoi(params["id"])
		if err == nil {
			for i := 0; i < num; i++ {
				w.Write([]byte("Hello, World!\n"))
				w.(http.Flusher).Flush()
				time.Sleep(1 * time.Second)
			}
		} else {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
		}
		// w.Write([]byte("Hello, World!\n"))
	}).Methods("GET")

	// Add middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)
}
