package twitlistserver

import (
	"errors"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/eric-fouillet/anaconda"
)

const refreshIntervalMin = 20 * time.Minute

// RealTwitterClient is a twitter client interacting with the Twitter API
type RealTwitterClient struct {
	api            *anaconda.TwitterApi
	lists          []anaconda.List
	lastUpdateTime time.Time
	authenticated  bool
}

// TwitterClient is an interface representing a Twitter client.
// Using an interface allows to mock the client and test more easily
type TwitterClient interface {
	authenticate() error
	close()
	getSelfID() (int64, error)
	GetListMembers(id int64) ([]anaconda.User, error)
	GetAllLists() ([]anaconda.List, error)
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
	v := url.Values{}
	v.Set("count", "30")
	return tc.api.GetListMembers(id, v)
}

// GetAllLists gets all lists for the authenticated user.
func (tc *RealTwitterClient) GetAllLists() ([]anaconda.List, error) {
	// Refresh the lists only every REFRESH_INTERVAL_MIN
	if tc.lists != nil && len(tc.lists) > 0 && time.Since(tc.lastUpdateTime) < refreshIntervalMin {
		log.Println("Re-use cached lists")
		return tc.lists, nil
	}
	id, err := tc.getSelfID()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("count", "30")
	tc.lists, err = tc.api.GetListsOwnedBy(id, v)
	tc.lastUpdateTime = time.Now()
	return tc.lists, nil
}
