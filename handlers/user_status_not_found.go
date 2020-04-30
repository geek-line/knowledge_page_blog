package handlers

import (
	"net/http"
	"text/template"

	"../structs"
)

// StatusNotFoundHandler に対するハンドラ
func StatusNotFoundHandler(w http.ResponseWriter, r *http.Request, auth bool) {
	header := newHeader(false)
	if auth {
		header.IsLogin = true
	}
	t := template.Must(template.ParseFiles("template/404.html", "template/_header.html", "template/_footer.html"))
	if err := t.Execute(w, struct {
		Header structs.Header
	}{
		Header: header,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
