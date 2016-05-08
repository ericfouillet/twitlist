package main

import (
	"net/http"
	"strconv"
)

func listHandler(w http.ResponseWriter, r *http.Request, tc *TwitterClient) {
	v := r.URL.Query()
	idStr := v.Get("id")
	if idStr == "" {
		http.Error(w, "No list ID specified", http.StatusInternalServerError)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Id has an incorrect format "+idStr, http.StatusInternalServerError)
	}
	users, err := tc.GetListMembers(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	renderTemplateUser("list", w, users)
}
