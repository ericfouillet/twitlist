package twitlistserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/eric-fouillet/anaconda"
)

// ListGet holds the result of a request to get a list of list members
type ListGet struct {
	ID      int64
	Members []anaconda.User
}

// ListHandler handles GET and POST to /lists/list/{id}
func ListHandler(w http.ResponseWriter, r *http.Request, tc TwitterClient) error {
	id, err := getListID(r)
	if err != nil {
		return fmt.Errorf("Id has an incorrect format %v", id)
	}
	switch r.Method {
	case "GET":
		return listHandlerGet(w, r, tc, id)
	case "POST":
		return listHandlerPost(w, r, tc, id)
	default:
		return errors.New(fmt.Sprintln("Unsupported method", r.Method))
	}
}

// listHandlerGet handles GET requests to /lists/list/{id}
func listHandlerGet(w http.ResponseWriter, r *http.Request, tc TwitterClient, listID int64) error {
	users, err := tc.GetListMembers(listID)
	if err != nil {
		return err
	}
	render := ListGet{listID, users}
	SetHeader(w, "GET")
	return json.NewEncoder(w).Encode(render)
}

// MemberIDs is a slice of members IDs (as int64)
type MemberIDs []struct{ ID int64 }

// listHandlerPut handles POST requests to /lists/list/{id}
func listHandlerPost(w http.ResponseWriter, r *http.Request, tc TwitterClient, listID int64) error {
	listID, err := getListID(r)
	if err != nil {
		return fmt.Errorf("Id has an incorrect format %v", listID)
	}
	var members MemberIDs
	if err := json.NewDecoder(r.Body).Decode(&members); err != nil {
		return err
	}
	membersList := make([]int64, 0, len(members))
	for _, m := range members {
		membersList = append(membersList, m.ID)
	}
	newMembers, err := tc.UpdateListMembers(listID, membersList)
	if err != nil {
		return err
	}
	render := ListGet{listID, newMembers}
	SetHeader(w, "PUT")
	return json.NewEncoder(w).Encode(render)
}

func getListID(r *http.Request) (int64, error) {
	rawPath := r.URL.Path
	last := strings.LastIndex(rawPath, "/")
	if last == -1 {
		return -1, errors.New("No list ID specified")
	}
	idStr := rawPath[last+1:]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return -1, errors.New("Id has an incorrect format " + idStr)
	}
	return id, nil
}
