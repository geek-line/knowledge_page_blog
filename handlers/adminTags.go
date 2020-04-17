package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

//AdminTagsHandler /admin/tagsに対するハンドラ
func AdminTagsHandler(w http.ResponseWriter, r *http.Request, env map[string]string) {
	session, _ := store.Get(r, "cookie-name")
	header := newHeader(false)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		header.IsLogin = true
	}
	db, err := sql.Open("mysql", env["SQL_ENV"])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	switch {
	case r.Method == "GET":
		rows, err := db.Query("SELECT id, name FROM tags")
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		var tags []Tag

		for rows.Next() {
			var tag Tag
			err := rows.Scan(&tag.ID, &tag.Name)
			if err != nil {
				panic(err.Error())
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
		_, err = db.Query("INSERT INTO tags(name, created_at, updated_at) VALUES(?, ?, ?)", name, createdAt, updatedAt)
		if err != nil {
			panic(err.Error())
		}
		http.Redirect(w, r, "/admin/tags/", http.StatusFound)
	case r.Method == "PUT":
		id, _ := strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("name")
		updatedAt := time.Now()
		_, err = db.Query("UPDATE tags SET name = ?, updated_at = ? WHERE id = ?", name, updatedAt, id)
		if err != nil {
			panic(err.Error())
		}
	case r.Method == "DELETE":
		id, _ := strconv.Atoi(r.FormValue("id"))
		_, err = db.Query("DELETE FROM tags WHERE id = ?", id)
		if err != nil {
			panic(err.Error())
		}
		_, err = db.Query("DELETE FROM knowledges_tags WHERE tag_id = ?", id)
		if err != nil {
			panic(err.Error())
		}
	}
}
