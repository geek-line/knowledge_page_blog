package handlers

import (
	"net/http"

	"../config"
	"../routes"
	"github.com/gorilla/sessions"
)

//AdminLogoutHandler /admin/logoutに対するハンドラ
func AdminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	store := sessions.NewCookieStore([]byte(config.SessionKey))
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, routes.AdminLoginPath, http.StatusFound)
}
