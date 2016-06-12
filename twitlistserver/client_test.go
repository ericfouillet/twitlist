package twitlistserver

import (
	"os"
	"testing"

	"github.com/eric-fouillet/anaconda"
)

func TestAuthenticate(t *testing.T) {
	params := []struct{ consKey, consSecret, accKey, accSecret, expected string }{
		{"", "abc", "def", "ghi", "error"},
		{"abc", "", "def", "ghi", "error"},
		{"abc", "abc", "", "ghi", "error"},
		{"ghi", "abc", "def", "", "error"},
	}
	tc := new(RealTwitterClient)
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

func TestUserDiff(t *testing.T) {
	params := []struct {
		existing  []anaconda.User
		requested int64arr
		added     []int64
		destroyed []int64
	}{
		{[]anaconda.User{anaconda.User{Id: 1, Name: "user1"}, anaconda.User{Id: 2, Name: "user2"}}, int64arr{1, 3}, []int64{3}, []int64{2}},
		{[]anaconda.User{}, int64arr{1, 3}, []int64{1, 3}, []int64{}},
		{[]anaconda.User{anaconda.User{Id: 1, Name: "user1"}}, int64arr{3}, []int64{3}, []int64{1}},
		{[]anaconda.User{anaconda.User{Id: 1, Name: "user1"}, anaconda.User{Id: 2, Name: "user2"}}, int64arr{1, 2}, []int64{}, []int64{}},
		{[]anaconda.User{anaconda.User{Id: 1, Name: "user1"}, anaconda.User{Id: 2, Name: "user2"}, anaconda.User{Id: 3, Name: "user3"}, anaconda.User{Id: 5, Name: "user5"}, anaconda.User{Id: 7, Name: "user7"}}, int64arr{1, 2, 5, 6}, []int64{6}, []int64{3, 7}},
	}
	for _, v := range params {
		added, destroyed := diffUsers(v.existing, v.requested)
		if !compareArrays(added, v.added) {
			t.Fail()
		}
		if !compareArrays(destroyed, v.destroyed) {
			t.Fail()
		}
	}
}

func compareArrays(a []int64, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
