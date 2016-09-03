package twitlistserver

import (
	"encoding/json"
	"net/http"

	"github.com/eric-fouillet/anaconda"
)

// TwitterLists holds a list of TwitterList.
// Those structs contain less information than what the Twitter API returns
// for simplicity.
type TwitterLists struct {
	TLists []TwitterList `json:"twitterLists"`
}

// ListsHandler handles GET requests to /lists
// Returns an array of lists
func ListsHandler(w http.ResponseWriter, r *http.Request, tc TwitterClient) error {
	lists, err := tc.GetAllLists()
	if err != nil {
		return err
	}
	res := makeLists(lists)
	SetHeader(w, "GET")
	return json.NewEncoder(w).Encode(res)
}

func makeLists(alists []anaconda.List) TwitterLists {
	var lists []TwitterList
	for _, ll := range alists {
		lists = append(lists, makeList(ll.Id, ll.Name, ll.Description, nil))
	}
	return TwitterLists{TLists: lists}
}
