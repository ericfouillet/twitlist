package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

var templates = template.Must(template.ParseFiles("lists.html", "home.html"))

func homeHandler(w http.ResponseWriter, r *http.Request, api *anaconda.TwitterApi) {
	//renderTemplate("home", w)
	fmt.Fprintf(w, "Hi")
}

func listsHandler(w http.ResponseWriter, r *http.Request, api *anaconda.TwitterApi) {
	lists, err := getAllLists(api)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	renderTemplate("lists", w, lists)
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, api *anaconda.TwitterApi), api *anaconda.TwitterApi) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, api)
	}
}

func main() {
	api := authenticate()
	http.HandleFunc("/", makeHandler(homeHandler, api))
	http.HandleFunc("/lists", makeHandler(listsHandler, api))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func authenticate() *anaconda.TwitterApi {
	consumerKey, consumerSecret := os.Getenv("TWIT_CONSUMER_KEY"), os.Getenv("TWIT_CONSUMER_SECRET")
	accessToken, accessTokenSecret := os.Getenv("TWIT_ACCESS_TOKEN"), os.Getenv("TWIT_ACCESS_TOKEN_SECRET")
	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessTokenSecret == "" {
		log.Fatal("Missing env variables")
	}
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	return anaconda.NewTwitterApi(accessToken, accessTokenSecret)
}

func getAllLists(api *anaconda.TwitterApi) ([]anaconda.List, error) {
	v := url.Values{}
	u, err := api.GetSelf(v)
	if err != nil {
		log.Fatal("Could not get current user")
	}
	v2 := url.Values{}
	v2.Set("count", "30")
	return api.GetListsOwnedBy(u.Id, v2)
}

func renderTemplate(tmpl string, w http.ResponseWriter, v []anaconda.List) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
