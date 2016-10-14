package twitlistserver

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/ericfouillet/anaconda"
)

const refreshIntervalMin = 20 * time.Minute

// RealTwitterClient is a twitter client interacting with the Twitter API
type RealTwitterClient struct {
	api            *anaconda.TwitterApi
	lists          []anaconda.List
	listMembers    map[int64][]anaconda.User
	friends        []anaconda.User
	lastUpdateTime time.Time
	authenticated  bool
}

func (tc *RealTwitterClient) authenticate() error {
	consumerKey, consumerSecret := os.Getenv("TWIT_CONSUMER_KEY"), os.Getenv("TWIT_CONSUMER_SECRET")
	accessToken, accessTokenSecret := os.Getenv("TWIT_ACCESS_TOKEN"), os.Getenv("TWIT_ACCESS_TOKEN_SECRET")
	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessTokenSecret == "" {
		return errors.New("Missing env variables")
	}
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	tc.api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	//tc.api.SetDelay(10 * time.Second)
	//tc.api.ReturnRateLimitError(true)
	return nil
}

func (tc *RealTwitterClient) close() {
	tc.api.Close()
}

func (tc *RealTwitterClient) getSelfID() (int64, error) {
	v := url.Values{}
	u, err := tc.api.GetSelf(v)
	if err != nil {
		return -1, err
	}
	return u.Id, nil
}

// GetListMembers retrieves all members of a list owned by the currently
// authenticated user.
func (tc *RealTwitterClient) GetListMembers(id int64) ([]anaconda.User, error) {
	if tc.listMembers == nil {
		tc.listMembers = make(map[int64][]anaconda.User)
	}
	// Refresh the list members only every REFRESH_INTERVAL_MIN
	members, ok := tc.listMembers[id]
	if ok && time.Since(tc.lastUpdateTime) < refreshIntervalMin {
		log.Println("Re-using cached list members")
		return members, nil
	}
	v := url.Values{}
	v.Set("count", "30")
	var err error
	tc.listMembers[id], err = tc.api.GetListMembers(id, v)
	if err != nil {
		return nil, err
	}
	tc.lastUpdateTime = time.Now()
	return tc.listMembers[id], nil
}

// GetAllLists gets all lists for the authenticated user.
func (tc *RealTwitterClient) GetAllLists() ([]anaconda.List, error) {
	// Refresh the lists only every REFRESH_INTERVAL_MIN
	if tc.lists != nil && len(tc.lists) > 0 && time.Since(tc.lastUpdateTime) < refreshIntervalMin {
		log.Println("Re-using cached lists")
		return tc.lists, nil
	}
	id, err := tc.getSelfID()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("count", "30")
	tc.lists, err = tc.api.GetListsOwnedBy(id, v)
	if err != nil {
		return nil, err
	}
	tc.lastUpdateTime = time.Now()
	return tc.lists, nil
}

func (tc *RealTwitterClient) GetAllFriends() ([]anaconda.User, error) {
	// Refresh the lists only every REFRESH_INTERVAL_MIN
	if tc.friends != nil && len(tc.friends) > 0 && time.Since(tc.lastUpdateTime) < refreshIntervalMin {
		log.Println("Re-using cached friends")
		return tc.friends, nil
	}
	v := url.Values{}
	v.Set("count", "30")
	friendsCursor, err := tc.api.GetFriendsList(v)
	if err != nil {
		return nil, err
	}
	tc.friends = friendsCursor.Users
	return tc.friends, nil
}

// UpdateListMembers updates the members of a list based on the
// members provided by the user (requestedMembers).
func (tc *RealTwitterClient) UpdateListMembers(listID int64, requestedMembers int64arr) ([]anaconda.User, error) {
	id, err := tc.getSelfID()
	if err != nil {
		return nil, err
	}

	existingMembers, err := tc.GetListMembers(listID)
	if err != nil {
		return nil, err
	}
	added, unchanged, destroyed := diffUsers(existingMembers, requestedMembers)

	var newUsers, destroyedUsers []anaconda.User

	if len(added) > 0 {
		v := url.Values{}
		v.Set("count", "30")
		newUsers, err = tc.api.AddUsersToList(id, added, v)
		if err != nil {
			return nil, err
		}
	}
	for _, nu := range newUsers {
		fmt.Printf("Added user %v\n", nu.Id)
	}

	if len(destroyed) > 0 {
		v2 := url.Values{}
		v2.Set("count", "30")
		destroyedUsers, err = tc.api.RemoveUsersFromList(id, destroyed, v2)
		if err != nil {
			return nil, err
		}
	}
	for _, du := range destroyedUsers {
		fmt.Printf("Removed user %v\n", du.Id)
	}

	tc.listMembers[listID] = updateMemberList(tc.listMembers[listID], unchanged, newUsers)

	return tc.listMembers[listID], nil
}
