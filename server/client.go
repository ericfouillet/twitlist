package server

import (
	"errors"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/eric-fouillet/anaconda"
)

const REFRESH_INTERVAL_MIN = 20 * time.Minute

type RealTwitterClient struct {
	api            *anaconda.TwitterApi
	lists          []anaconda.List
	lastUpdateTime time.Time
	authenticated  bool
}

type TwitterClient interface {
	authenticate() error
	close()
	getSelfId() (int64, error)
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

func (tc *RealTwitterClient) getSelfId() (int64, error) {
	v := url.Values{}
	u, err := tc.api.GetSelf(v)
	if err != nil {
		return -1, err
	}
	return u.Id, nil
}

func (tc *RealTwitterClient) GetListMembers(id int64) ([]anaconda.User, error) {
	v := url.Values{}
	v.Set("count", "30")
	return tc.api.GetListMembers(id, v)
}

func (tc *RealTwitterClient) GetAllLists() ([]anaconda.List, error) {
	// Refresh the lists only every REFRESH_INTERVAL_MIN
	if tc.lists != nil && len(tc.lists) > 0 && time.Since(tc.lastUpdateTime) < REFRESH_INTERVAL_MIN {
		log.Println("Re-use cached lists")
		return tc.lists, nil
	}
	id, err := tc.getSelfId()
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("count", "30")
	tc.lists, err = tc.api.GetListsOwnedBy(id, v)
	tc.lastUpdateTime = time.Now()
	return tc.lists, nil
}
