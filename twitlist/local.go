package main

import (
	"log"
	"net/http"

	"github.com/eric-fouillet/twitlist/twitlistserver"
)

func main() {
	twitlistserver.RegisterHandlers()
	http.Handle("/", http.FileServer(http.Dir("static")))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
