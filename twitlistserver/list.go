package twitlistserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/eric-fouillet/anaconda"
)

type listGet struct {
	ID      int64
	Members []anaconda.User
}

// listHandler handles GET and POST to /lists/list/{id}
func listHandler(w http.ResponseWriter, r *http.Request, tc TwitterClient) error {
	log.Println("Entered listHandler", r.Method)
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
	render := listGet{listID, users}
	return json.NewEncoder(w).Encode(render)
}

type memberIDs []struct{ ID int64 }

// listHandlerPut handles POST requests to /lists/list/{id}
func listHandlerPost(w http.ResponseWriter, r *http.Request, tc TwitterClient, listID int64) error {
	listID, err := getListID(r)
	if err != nil {
		return fmt.Errorf("Id has an incorrect format %v", listID)
	}
	var members memberIDs
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
	render := listGet{listID, newMembers}
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
