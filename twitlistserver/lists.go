package twitlistserver

import (
	"encoding/json"
	"net/http"

	"github.com/eric-fouillet/anaconda"
)

// ListsHandler handles GET requests to /lists
// Returns an array of lists
func ListsHandler(w http.ResponseWriter, r *http.Request, tc TwitterClient) error {
	lists, err := tc.GetAllLists()
	if err != nil {
		return err
	}
	res := struct{ Lists []anaconda.List }{lists}
	SetHeader(w, "GET")
	return json.NewEncoder(w).Encode(res)
}
