package main

import (
	"html/template"
	"log"
	"net/http"

	//"github.com/ChimeraCoder/anaconda"
	"github.com/eric-fouillet/anaconda"
)

var templates = template.Must(template.ParseFiles("lists.html", "list.html"))

func main() {
	tc := new(TwitterClient)
	err := tc.authenticate()
	defer tc.close()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", makeHandler(listsHandler, tc))
	http.HandleFunc("/list", makeHandler(listHandler, tc))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func makeHandler(fn func(w http.ResponseWriter,
	r *http.Request,
	tc *TwitterClient),
	tc *TwitterClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, tc)
	}
}

func renderTemplateList(tmpl string, w http.ResponseWriter, v []anaconda.List) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateUser(tmpl string, w http.ResponseWriter, v []anaconda.User) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
