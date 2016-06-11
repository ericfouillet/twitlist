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
	http.HandleFunc(pathPrefix+"list/", makeHandler(listHandler, tc))
}

func makeHandler(fn func(w http.ResponseWriter,
	r *http.Request,
	tc TwitterClient) error,
	tc TwitterClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r, tc)
		if err == nil {
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
