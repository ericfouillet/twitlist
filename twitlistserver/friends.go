package twitlistserver

import (
	"encoding/json"
	"net/http"
)

// TwitterUser holds limited details about a Twitter list member.
// Only essential details are used, since the UI does not need to have
// everything to maintain the lists.
type TwitterUser struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TwitterFriends is
type TwitterFriends struct {
	Friends []TwitterUser `json:"friends"`
}

// FriendsHandler handles GET requests on /users.
// It returns all contacts of the current user on Twitter.
// This allows the UI to offer a list of users to add to a list.
func FriendsHandler(w http.ResponseWriter, r *http.Request, tc TwitterClient) error {
	users, err := tc.GetAllFriends()
	if err != nil {
		return err
	}
	allUsers := make([]TwitterUser, 0)
	for _, u := range users {
		allUsers = append(allUsers, TwitterUser{ID: u.Id, Name: u.Name, Description: u.Description})
	}
	render := TwitterFriends{Friends: allUsers}
	SetHeader(w, "GET")
	return json.NewEncoder(w).Encode(render)
}
