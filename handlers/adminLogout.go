package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"
)

//AdminLogoutHandler /admin/logoutに対するハンドラ
func AdminLogoutHandler(w http.ResponseWriter, r *http.Request, env map[string]string, db *sql.DB) {
	store := sessions.NewCookieStore([]byte(env["SESSION_KEY"]))
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/admin/login/", http.StatusFound)
}
