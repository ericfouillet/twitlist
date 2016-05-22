package twitlistserver

import (
	"log"
	"net/http"
)

const pathPrefix = "/lists/"

// RegisterHandlers registers all handlers to serve requests.
func RegisterHandlers() {
	tc := new(DummyTwitterClient)
	err := tc.authenticate()
	defer tc.close()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc(pathPrefix, makeHandler(listsHandler, tc))
	http.HandleFunc(pathPrefix+"list", makeHandler(listHandler, tc))
}

func makeHandler(fn func(w http.ResponseWriter,
	r *http.Request,
	tc TwitterClient),
	tc TwitterClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, tc)
	}
}
