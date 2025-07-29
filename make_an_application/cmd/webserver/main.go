package main

import (
	"log"
	"net/http"

	poker "github.com/espennoreng/learn-go-with-tests/make_an_application"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()
	
	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(http.ListenAndServe(":8080", server))
	}
}
