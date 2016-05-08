package main

import "net/http"

func listsHandler(w http.ResponseWriter, r *http.Request, tc *TwitterClient) {
	err := tc.GetAllLists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	renderTemplateList("lists", w, tc.lists)
}
