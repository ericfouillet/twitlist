package twitlistserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetFriends(t *testing.T) {
	// Prepare the client
	c := new(DummyTwitterClient)
	c.authenticate()

	// Prepare the request
	req, err := http.NewRequest("GET", "/friends", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare the handler function
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MakeHandler(FriendsHandler, c))
	handler.ServeHTTP(rr, req)

	if !strings.Contains(rr.Body.String(), `"name":"1user9"`) {
		t.Fail()
	}
	result := new(TwitterFriends)
	if err := json.NewDecoder(rr.Body).Decode(result); err != nil {
		t.Fail()
	}
	if len(result.Friends) != 50 {
		t.Fail()
	}
}
