package main

import (
	"log"
	"net/http"

	dependencyinjection "github.com/espennoreng/learn-go-with-tests/dependency_injection"
)

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
	dependencyinjection.Greet(w, "Espen")
}

func main() {
	log.Fatal(http.ListenAndServe(":5001", http.HandlerFunc(MyGreeterHandler)))
}
