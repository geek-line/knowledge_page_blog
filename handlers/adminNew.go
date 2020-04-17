package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
)

func newHeader(isLogin bool) Header {
	return Header{IsLogin: isLogin}
}

//AdminNewHandler /admin/newに対するハンドラ
func AdminNewHandler(w http.ResponseWriter, r *http.Request, env map[string]string) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
	}

	db, err := sql.Open("mysql", env["SQL_ENV"])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, name FROM tags")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	var tags []Tag

	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.Id, &tag.Name)
		if err != nil {
			panic(err.Error())
		}
		tags = append(tags, tag)
	}
	t := template.Must(template.ParseFiles("template/admin_new.html", "template/_header.html"))
	header := newHeader(true)
	if err := t.Execute(w, struct {
		Header Header
		Tags   []Tag
	}{
		Header: header,
		Tags:   tags,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
