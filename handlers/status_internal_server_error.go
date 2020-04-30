package handlers

import (
	"html/template"
	"net/http"
)

// StatusInternalServerError に対するハンドラ
func StatusInternalServerError(w http.ResponseWriter, r *http.Request, auth bool) {
	header := newHeader(false)
	if auth {
		header.IsLogin = true
	}
	t := template.Must(template.ParseFiles("template/500.html", "template/_header.html", "template/_footer.html"))
	t.Execute(w, struct {
		Header Header
	}{
		Header: header,
	})
}
