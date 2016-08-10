package twitlistserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eric-fouillet/anaconda"
)

func TestGetLists(t *testing.T) {
	// Prepare the client
	c := new(DummyTwitterClient)
	c.authenticate()

	// Prepare the request
	req, err := http.NewRequest("GET", "/lists", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare the handler function
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MakeHandler(ListsHandler, c))
	handler.ServeHTTP(rr, req)

	if !strings.Contains(rr.Body.String(), `"name":"list1"`) {
		t.Fail()
	}
	result := new(struct{ Lists []anaconda.List })
	if err := json.NewDecoder(rr.Body).Decode(result); err != nil {
		t.Fail()
	}
	if len(result.Lists) != 5 {
		t.Fail()
	}
}
