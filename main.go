package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/david-galdamez/search-engine/database"
	"github.com/david-galdamez/search-engine/handlers"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error getting .env file: %v", err.Error())
	}

	err = database.CreateBuckets()
	if err != nil {
		log.Fatalf("error creating buckets")
	}

	router := http.NewServeMux()

	router.HandleFunc("POST /index", handlers.Index)
	router.HandleFunc("GET /search", handlers.Search)
	router.HandleFunc("GET /docs/{id}", handlers.Docs)
	router.HandleFunc("GET /counter", handlers.Counter)

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Health Check"))
	})

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatalf("PORT env not found")
	}

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}

	fmt.Print("Server listening on port: 8080\n")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("error listening to server: %v\n", err)
	}

}
