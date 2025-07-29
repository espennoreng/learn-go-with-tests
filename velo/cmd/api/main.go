package main

import (
	"log"
	"net/http"

	"github.com/espennoreng/learn-go-with-tests/velo"
)

func main() {
	store := velo.NewInMemoryAppStore()

	server := velo.NewAppServer(store)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
