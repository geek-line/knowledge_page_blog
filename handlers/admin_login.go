package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"../routes"

	"github.com/gorilla/sessions"
)

//AdminLoginHandler /admin/loginに対するハンドラ
func AdminLoginHandler(w http.ResponseWriter, r *http.Request, env map[string]string, db *sql.DB) {
	store := sessions.NewCookieStore([]byte(env["SESSION_KEY"]))
	if r.Method == "GET" {
		session, _ := store.Get(r, "cookie-name")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			t := template.Must(template.ParseFiles("template/admin_login.html", "template/_header.html"))
			header := newHeader(false)
			if err := t.Execute(w, struct {
				Header Header
			}{
				Header: header,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			session.Values["authenticated"] = true
			session.Save(r, w)
			http.Redirect(w, r, routes.AdminKnowledgesPath, http.StatusFound)
		}
	} else {
		email := r.FormValue("email")
		password := r.FormValue("password")
		var correctPassword string
		if err := db.QueryRow("SELECT password FROM admin_user WHERE email = ?", email).Scan(&correctPassword); err != nil {
			http.Redirect(w, r, routes.AdminLoginPath, http.StatusFound)
		}
		if correctPassword == password {
			session, _ := store.Get(r, "cookie-name")
			session.Values["authenticated"] = true
			session.Save(r, w)
			http.Redirect(w, r, routes.AdminKnowledgesPath, http.StatusFound)

		} else {
			http.Redirect(w, r, routes.AdminLoginPath, http.StatusFound)
			return
		}
	}
}
