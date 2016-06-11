package twitlistserver

import (
	"encoding/json"
	"net/http"

	"github.com/eric-fouillet/anaconda"
)

// listsHandler handles GET requests to /lists
// Returns an array of lists
func listsHandler(w http.ResponseWriter, r *http.Request, tc TwitterClient) error {
	lists, err := tc.GetAllLists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res := struct{ Lists []anaconda.List }{lists}
	return json.NewEncoder(w).Encode(res)
}
