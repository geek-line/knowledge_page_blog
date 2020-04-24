package handlers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

// StatusInternalServerError に対するハンドラ
func StatusInternalServerError(w http.ResponseWriter, r *http.Request, env map[string]string) {
	store := sessions.NewCookieStore([]byte(env["SESSION_KEY"]))
	session, _ := store.Get(r, "cookie-name")
	header := newHeader(false)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		header.IsLogin = true
	}
	t := template.Must(template.ParseFiles("template/500.html", "template/_header.html", "template/_footer.html"))
	if err := t.Execute(w, struct {
		Header Header
	}{
		Header: header,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
