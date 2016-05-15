package server

import (
	"html/template"
	"net/http"

	"github.com/eric-fouillet/anaconda"
)

var listsTempl = template.Must(template.New("list").Parse(listsTemplateHtml))

func listsHandler(w http.ResponseWriter, r *http.Request, tc TwitterClient) {
	lists, err := tc.GetAllLists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	renderTemplateList(w, lists)
}

func renderTemplateList(w http.ResponseWriter, v []anaconda.List) {
	if err := listsTempl.Execute(w, v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const listsTemplateHtml = `
<h1>Lists</h1>
<ul>
{{range $entry := .}}
<li><a href="/lists/list?id={{$entry.Id}}">{{$entry.Name}}</a></li>
{{end}}
</ul>
`
