package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func newHeader(isLogin bool) Header {
	return Header{IsLogin: isLogin}
}

//AdminNewHandler /admin/newに対するハンドラ
func AdminNewHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT id, name FROM tags")
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer rows.Close()
	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			log.Print(err.Error())
			return
		}
		tags = append(tags, tag)
	}
	rows, err = db.Query("SELECT name, src FROM eyecatches")
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer rows.Close()
	var eyecatches []EyeCatch
	for rows.Next() {
		var eyecatch EyeCatch
		err := rows.Scan(&eyecatch.Name, &eyecatch.Src)
		if err != nil {
			log.Print(err.Error())
			return
		}
		eyecatches = append(eyecatches, eyecatch)
	}
	t := template.Must(template.ParseFiles("template/admin_new.html", "template/_header.html"))
	header := newHeader(true)
	if err := t.Execute(w, struct {
		Header     Header
		Tags       []Tag
		EyeCatches []EyeCatch
	}{
		Header:     header,
		Tags:       tags,
		EyeCatches: eyecatches,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
