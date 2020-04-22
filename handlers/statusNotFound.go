package handlers

import (
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

// StatusNotFoundHandler に対するハンドラ
func StatusNotFoundHandler(w http.ResponseWriter, r *http.Request, env map[string]string) {
	store := sessions.NewCookieStore([]byte(env["SESSION_KEY"]))
	session, _ := store.Get(r, "cookie-name")
	header := newHeader(false)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		header.IsLogin = true
	}
	t := template.Must(template.ParseFiles("template/404.html", "template/_header.html", "template/_footer.html"))
	if err := t.Execute(w, struct {
		Header Header
	}{
		Header: header,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
