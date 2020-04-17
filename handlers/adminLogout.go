package handlers

import "net/http"

//AdminLogoutHandler /admin/logoutに対するハンドラ
func AdminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/admin/login/", http.StatusFound)
}
