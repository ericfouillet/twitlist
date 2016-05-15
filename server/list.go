package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/eric-fouillet/anaconda"
)

var listTempl = template.Must(template.New("list").Parse(listTemplateHtml))

type listGet struct {
	Id      int64
	Members []anaconda.User
}

func listHandler(w http.ResponseWriter, r *http.Request, tc TwitterClient) {
	switch r.Method {
	case "GET":
		listHandlerGet(w, r, tc)
	case "PUT":
		listHandlerPut(w, r, tc)
	}
}

func listHandlerGet(w http.ResponseWriter, r *http.Request, tc TwitterClient) {
	v := r.URL.Query()
	idStr := v.Get("id")
	if idStr == "" {
		http.Error(w, "No list ID specified", http.StatusInternalServerError)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Id has an incorrect format "+idStr, http.StatusInternalServerError)
		return
	}
	users, err := tc.GetListMembers(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render := listGet{id, users}
	renderTemplateUser("list", w, render)
}

func listHandlerPut(w http.ResponseWriter, r *http.Request, tc TwitterClient) {
	v := r.URL.Query()
	listIdStr := v.Get("listId")
	memberIdsStr := v.Get("memberIds")
	if listIdStr == "" || memberIdsStr == "" {
		http.Error(w, "No list ID or member Ids specified", http.StatusInternalServerError)
		return
	}
	listId, err := strconv.ParseInt(listIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Id has an incorrect format "+listIdStr, http.StatusInternalServerError)
		return
	}
	membersIdsList := strings.Split(memberIdsStr, ",")
	memberIds := make([]int64, 0, len(membersIdsList))
	for _, mIdStr := range membersIdsList {
		mId, err := strconv.ParseInt(mIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Id has an incorrect format "+mIdStr, http.StatusInternalServerError)
			return
		}
		memberIds = append(memberIds, mId)
	}
	fmt.Fprintln(w, listId, memberIds)
}

func renderTemplateUser(tmpl string, w http.ResponseWriter, v listGet) {
	if err := listTempl.Execute(w, v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const listTemplateHtml = `
<h1>List detail</h1>
<form action="/list/{{.Id}}">
<table border="1">
  <!--<th>
    <td>Name</td>
    <td>In list</td>
  </th>-->
{{range $entry := .Members}}
<tr>
<td>{{$entry.Name }}</td>
<td><input name="{{$entry.Id}}" type="checkbox" checked/></td>
</tr>
{{end}}
</table>
<input type="submit" value="Save"/>
</form>
`
