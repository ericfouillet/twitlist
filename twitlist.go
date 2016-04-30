package main

import (
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func authenticate() *anaconda.TwitterApi {
	consumerKey, consumerSecret := os.Getenv("TWIT_CONSUMER_KEY"), os.Getenv("TWIT_CONSUMER_SECRET")
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	accessToken, accessTokenSecret := os.Getenv("TWIT_ACCESS_TOKEN"), os.Getenv("TWIT_ACCESS_TOKEN_SECRET")
	return anaconda.NewTwitterApi(accessToken, accessTokenSecret)
}

func getAllLists(api *anaconda.TwitterApi) (anaconda.List, err) {
	v := url.Values{}
	u, err := api.GetSelf(v)
	if err != nil {
		log.Fatal("Could not get current user")
	}
	v := url.Values{}
	v.Set("count", "30")
	return api.GetListsOwnedBy(u.Id, value)
}
