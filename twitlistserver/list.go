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
	switch r.Method {
	case "GET":
		return listHandlerGet(w, r, tc)
	case "POST":
		return listHandlerPost(w, r, tc)
	}
	return errors.New(fmt.Sprintln("Unsupported method", r.Method))
}

// listHandler handles GET requests to /lists/list/{id}
func listHandlerGet(w http.ResponseWriter, r *http.Request, tc TwitterClient) error {
	id, err := getListID(r)
	if err != nil {
		return fmt.Errorf("Id has an incorrect format %v", id)
	}
	users, err := tc.GetListMembers(id)
	if err != nil {
		return err
	}
	render := listGet{id, users}
	return json.NewEncoder(w).Encode(render)
}

type memberIds []struct{ ID int64 }

// listHanlderPut handles POST requests to /lists/list/{id}
func listHandlerPost(w http.ResponseWriter, r *http.Request, tc TwitterClient) error {
	id, err := getListID(r)
	if err != nil {
		return fmt.Errorf("Id has an incorrect format %v", id)
	}
	var members memberIds
	if err := json.NewDecoder(r.Body).Decode(&members); err != nil {
		return err
	}
	fmt.Fprintln(w, id, members)
	return nil
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
