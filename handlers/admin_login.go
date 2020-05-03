package handlers

import (
	"html/template"
	"log"
	"net/http"

	"../config"
	"../models"
	"../routes"
	"../structs"
	"github.com/gorilla/sessions"
)

//AdminLoginHandler /admin/loginに対するハンドラ
func AdminLoginHandler(w http.ResponseWriter, r *http.Request, auth bool) {
	store := sessions.NewCookieStore([]byte(config.SessionKey))
	if r.Method == "GET" {
		if !auth {
			t := template.Must(template.ParseFiles("template/admin_login.html", "template/_header.html"))
			header := newHeader(false)
			if err := t.Execute(w, struct {
				Header structs.Header
			}{
				Header: header,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			session, _ := store.Get(r, "cookie-name")
			session.Values["authenticated"] = true
			session.Save(r, w)
			http.Redirect(w, r, routes.AdminKnowledgesPath, http.StatusFound)
		}
	} else {
		email := r.FormValue("email")
		password := r.FormValue("password")
		var correctPassword string
		correctPassword, err := models.GetPasswordFromEmail(email)
		if err != nil {
			log.Print(err.Error())
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
