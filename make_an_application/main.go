package main

import (
	"log"
	"net/http"
)

func main() {
	server := &PlayerServer{store: NewInMemoryPlayerStore()}
	log.Fatal(http.ListenAndServe(":8080", server))
}
