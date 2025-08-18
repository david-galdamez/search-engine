package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/david-galdamez/search-engine/database"
	"github.com/david-galdamez/search-engine/handlers"
)

func main() {

	db, err := database.GetDB()
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists([]byte("terms"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("docs"))
		if err != nil {
			return err
		}

		b, err := tx.CreateBucketIfNotExists([]byte("meta"))
		if err != nil {
			return err
		}

		err = b.Put([]byte("N"), []byte("0"))
		if err != nil {
			return err
		}

		return nil
	})

	db.Close()

	router := http.NewServeMux()
	router.HandleFunc("POST /index", handlers.Index)

	router.HandleFunc("GET /search", handlers.Search)

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Health Check"))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error listening to server: %v\n", err)
	}

	fmt.Print("Server listening on port: 8080")

}
