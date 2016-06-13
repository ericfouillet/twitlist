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

	updateMemberList(existingMembers, unchanged, newUsers)

	// TODO update members list
	return tc.GetListMembers(listID)
}
