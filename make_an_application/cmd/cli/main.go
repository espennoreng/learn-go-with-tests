package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/espennoreng/learn-go-with-tests/make_an_application"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer close()

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.StdOutAlerter), store)

	poker.NewCLI(os.Stdin, os.Stdout, game).PlayPoker()

}
