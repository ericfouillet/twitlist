package twitlistserver

import (
	"flag"
	"log"
	"net/http"

	"github.com/eric-fouillet/anaconda"
)

const pathPrefix = "/lists/"

var clientType = flag.String("client", "dummy", "The client type to use (real/dummy). Default is dummy.")

// TwitterClient is an interface representing a Twitter client.
// Using an interface allows to mock the client and test more easily
type TwitterClient interface {
	authenticate() error
	close()
	getSelfID() (int64, error)
	GetAllLists() ([]anaconda.List, error)
	GetListMembers(id int64) ([]anaconda.User, error)
	UpdateListMembers(listID int64, requestedMembers int64arr) ([]anaconda.User, error)
}

// RegisterHandlers registers all handlers to serve requests.
func RegisterHandlers() {
	var tc TwitterClient
	switch *clientType {
	case "real":
		tc = new(DummyTwitterClient)
	case "dummy":
		tc = new(RealTwitterClient)
	default:
		log.Fatal("Unsupported client type : " + *clientType)
	}
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
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}
