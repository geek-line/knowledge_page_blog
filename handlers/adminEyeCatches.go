package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

//AdminEyeCatchesHandler /admin/eyecatchesに対するハンドラ
func AdminEyeCatchesHandler(w http.ResponseWriter, r *http.Request, env map[string]string, db *sql.DB) {
	store := sessions.NewCookieStore([]byte(env["SESSION_KEY"]))
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
		return
	}
	header := newHeader(false)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		header.IsLogin = true
	}
	switch {
	case r.Method == "GET":
		rows, err := db.Query("SELECT id, name, src FROM eyecatches")
		if err != nil {
			StatusNotFoundHandler(w, r, env)
		}
		defer rows.Close()
		var eyecatches []EyeCatch
		for rows.Next() {
			var eyecatch EyeCatch
			if err := rows.Scan(&eyecatch.ID, &eyecatch.Name, &eyecatch.Src); err != nil {
				StatusNotFoundHandler(w, r, env)
			}
			eyecatches = append(eyecatches, eyecatch)
		}
		t := template.Must(template.ParseFiles("template/admin_eyecatches.html", "template/_header.html"))
		if err := t.Execute(w, struct {
			Header     Header
			EyeCatches []EyeCatch
		}{
			Header:     header,
			EyeCatches: eyecatches,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case r.Method == "POST":
		name := r.FormValue("name")
		src := r.FormValue("src")
		rows, err := db.Query("INSERT INTO eyecatches(name, src) VALUES(?, ?)", name, src)
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		}
		defer rows.Close()
		http.Redirect(w, r, "/admin/eyecatches//", http.StatusFound)
	case r.Method == "PUT":
		id, _ := strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("name")
		src := r.FormValue("src")
		rows, err := db.Query("UPDATE eyecatches SET name = ?, src = ? WHERE id = ?", name, src, id)
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		}
		defer rows.Close()
	case r.Method == "DELETE":
		id, _ := strconv.Atoi(r.FormValue("id"))
		rows, err := db.Query("DELETE FROM eyecatches WHERE id = ?", id)
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		}
		defer rows.Close()
	}
}
