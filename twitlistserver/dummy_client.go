package twitlistserver

import (
	"fmt"
	"time"

	"github.com/eric-fouillet/anaconda"
)

// DummyTwitterClient is a Twitter client serving fake data.
// For test only.
type DummyTwitterClient struct {
	api            *anaconda.TwitterApi
	lists          []anaconda.List
	listMembers    map[int64][]anaconda.User
	friends        []anaconda.User
	lastUpdateTime time.Time
	authenticated  bool
}

func (tc *DummyTwitterClient) authenticate() error {
	return nil
}

func (tc *DummyTwitterClient) close() {
}

func (tc *DummyTwitterClient) getSelfID() (int64, error) {
	return 1, nil
}

// GetListMembers retrieves all members of a list owned by the currently
// authenticated user.
func (tc *DummyTwitterClient) GetListMembers(id int64) ([]anaconda.User, error) {
	if tc.listMembers == nil {
		tc.listMembers = make(map[int64][]anaconda.User)
	}
	users := make([]anaconda.User, 0, 10)
	var add func()
	var idu int64 = 1
	add = func() {
		users = append(users, anaconda.User{Id: idu, Name: fmt.Sprint(id, "user", idu)})
		idu++
	}
	for i := 0; i < 10; i++ {
		add()
	}
	tc.listMembers[id] = users
	return users, nil
}

// GetAllLists gets all lists for the authenticated user.
func (tc *DummyTwitterClient) GetAllLists() ([]anaconda.List, error) {
	lists := make([]anaconda.List, 0, 5)
	var add func()
	var id int64 = 1
	add = func() {
		lists = append(lists, anaconda.List{Id: id, Name: fmt.Sprint("list", id)})
		id++
	}
	for i := 0; i < 5; i++ {
		add()
	}
	tc.lists = lists
	return tc.lists, nil
}

// GetAllFriends gets all friends for the authenticated user.
func (tc *DummyTwitterClient) GetAllFriends() ([]anaconda.User, error) {
	friends := make([]anaconda.User, 0, 50)
	var add func(lid int64)
	var uid int64 = 1
	add = func(lid int64) {
		friends = append(friends, anaconda.User{Id: uid, Name: fmt.Sprint(lid, "user", uid)})
		uid++
	}
	for i := 1; i <= 5; i++ {
		for j := 0; j < 10; j++ {
			add(int64(i))
		}
	}
	tc.friends = friends
	return tc.friends, nil
}

// UpdateListMembers updates the members of a list based on the new
// members received from the user.
func (tc *DummyTwitterClient) UpdateListMembers(listID int64, requestedMembers int64arr) ([]anaconda.User, error) {
	existingMembers, err := tc.GetListMembers(listID)
	if err != nil {
		return nil, err
	}
	added, unchanged, destroyed := diffUsers(existingMembers, requestedMembers)

	var newUsers []anaconda.User

	if len(added) > 0 {
		for _, nu := range added {
			fmt.Printf("Added user %v\n", nu)
			newUsers = append(newUsers, anaconda.User{Id: nu, Name: fmt.Sprint(listID, "user", nu)})
		}
	}

	if len(destroyed) > 0 {
		for _, du := range destroyed {
			fmt.Printf("Removed user %v\n", du)
		}
	}

	tc.listMembers[listID] = updateMemberList(existingMembers, unchanged, newUsers)

	return tc.listMembers[listID], nil
}
