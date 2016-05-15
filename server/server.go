package server

import (
	"log"
	"net/http"
)

const PathPrefix = "/lists/"

func RegisterHandlers() {
	tc := new(DummyTwitterClient)
	err := tc.authenticate()
	defer tc.close()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc(PathPrefix, makeHandler(listsHandler, tc))
	http.HandleFunc(PathPrefix+"list", makeHandler(listHandler, tc))
}

func makeHandler(fn func(w http.ResponseWriter,
	r *http.Request,
	tc TwitterClient),
	tc TwitterClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, tc)
	}
}
