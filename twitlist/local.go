package main

import (
	"log"
	"net/http"

	"github.com/eric-fouillet/twitlist/server"
)

func main() {
	server.RegisterHandlers()
	//http.Handle("/", http.FileServer(http.Dir("static")))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
