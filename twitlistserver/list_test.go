package twitlistserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetList(t *testing.T) {
	// Prepare the client
	c := new(DummyTwitterClient)
	c.authenticate()

	// Prepare the request
	req, err := http.NewRequest("GET", "/list/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare the handler function
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MakeHandler(ListHandler, c))
	handler.ServeHTTP(rr, req)

	if !strings.Contains(rr.Body.String(), `"id":1`) {
		t.Fail()
	}
	result := new(TwitterList)
	if err := json.NewDecoder(rr.Body).Decode(result); err != nil {
		t.Fail()
	}
	if result.ID != 1 {
		t.Fail()
	}
	if len(result.Members) != 10 {
		t.Fail()
	}
}

func TestPutList(t *testing.T) {
	// Prepare the client
	c := new(DummyTwitterClient)
	c.authenticate()

	// Prepare the request
	members := MemberIDs{{ID: 1}, {ID: 2}}
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(&members); err != nil {
		t.Fatal("Could not encode request")
	}
	req, err := http.NewRequest("POST", "/list/1", b)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare the handler function
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MakeHandler(ListHandler, c))
	handler.ServeHTTP(rr, req)

	if !strings.Contains(rr.Body.String(), `"id":1`) {
		t.Fail()
	}
	result := new(TwitterList)
	if err := json.NewDecoder(rr.Body).Decode(result); err != nil {
		t.Fail()
	}
	if result.ID != 1 {
		t.Fail()
	}
	if len(result.Members) != 2 {
		t.Fail()
	}
}
