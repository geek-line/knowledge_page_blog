package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"../routes"

	"github.com/gorilla/sessions"
)

//AdminTagsHandler /admin/tagsに対するハンドラ
func AdminTagsHandler(w http.ResponseWriter, r *http.Request, env map[string]string, db *sql.DB) {
	store := sessions.NewCookieStore([]byte(env["SESSION_KEY"]))
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, routes.AdminLoginPath, http.StatusFound)
		return
	}
	header := newHeader(false)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		header.IsLogin = true
	}
	switch {
	case r.Method == "GET":
		rows, err := db.Query("SELECT id, name FROM tags")
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		}
		defer rows.Close()
		var tags []Tag

		for rows.Next() {
			var tag Tag
			err := rows.Scan(&tag.ID, &tag.Name)
			if err != nil {
				log.Print(err.Error())
				StatusInternalServerError(w, r, env)
				return
			}
			tags = append(tags, tag)
		}
		t := template.Must(template.ParseFiles("template/admin_tags.html", "template/_header.html"))
		if err := t.Execute(w, struct {
			Header Header
			Tags   []Tag
		}{
			Header: header,
			Tags:   tags,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case r.Method == "POST":
		name := r.FormValue("name")
		createdAt := time.Now()
		updatedAt := time.Now()
		rows, err := db.Query("INSERT INTO tags(name, created_at, updated_at) VALUES(?, ?, ?)", name, createdAt, updatedAt)
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		}
		defer rows.Close()
		http.Redirect(w, r, routes.AdminTagsPath, http.StatusFound)
	case r.Method == "PUT":
		id, _ := strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("name")
		updatedAt := time.Now()
		rows, err := db.Query("UPDATE tags SET name = ?, updated_at = ? WHERE id = ?", name, updatedAt, id)
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		}
		defer rows.Close()
	case r.Method == "DELETE":
		id, _ := strconv.Atoi(r.FormValue("id"))
		rows, err := db.Query("DELETE FROM tags WHERE id = ?", id)
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		}
		defer rows.Close()
		rows, err = db.Query("DELETE FROM knowledges_tags WHERE tag_id = ?", id)
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		}
		defer rows.Close()
	}
}