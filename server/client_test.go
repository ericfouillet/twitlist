package server

import (
	"os"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	params := []struct{ consKey, consSecret, accKey, accSecret, expected string }{
		{"", "abc", "def", "ghi", "error"},
		{"abc", "", "def", "ghi", "error"},
		{"", "abc", "def", "ghi", "error"},
		{"", "abc", "def", "ghi", "error"},
		{"", "abc", "def", "ghi", "error"},
		{"", "abc", "def", "ghi", "error"},
	}
	tc := new(TwitterClient)
	for _, v := range params {
		os.Setenv("TWIT_CONSUMER_KEY", v.consKey)
		os.Setenv("TWIT_CONSUMER_SECRET", v.consSecret)
		os.Setenv("TWIT_ACCESS_TOKEN", v.accKey)
		os.Setenv("TWIT_ACCESS_TOKEN_SECRET", v.accSecret)
		err := tc.authenticate()
		if v.expected != "" && err == nil {
			t.Fail()
		}
		if v.expected == "" && err != nil {
			t.Fail()
		}
		if err != nil && tc.api != nil {
			t.Fail()
		}
	}
}
