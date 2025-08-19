package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/david-galdamez/search-engine/handlers"
)

func main() {

	router := http.NewServeMux()

	router.HandleFunc("POST /index", handlers.Index)
	router.HandleFunc("GET /search", handlers.Search)
	router.HandleFunc("GET /docs/{id}", handlers.Docs)

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Health Check"))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error listening to server: %v\n", err)
	}

	fmt.Print("Server listening on port: 8080")

}
