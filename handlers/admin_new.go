package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"../routes"

	"github.com/gorilla/sessions"
)

func newHeader(isLogin bool) Header {
	return Header{IsLogin: isLogin}
}

//AdminNewHandler /admin/newに対するハンドラ
func AdminNewHandler(w http.ResponseWriter, r *http.Request, env map[string]string, db *sql.DB) {
	store := sessions.NewCookieStore([]byte(env["SESSION_KEY"]))
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, routes.AdminLoginPath, http.StatusFound)
		return
	}
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
	rows, err = db.Query("SELECT name, src FROM eyecatches")
	if err != nil {
		log.Print(err.Error())
		StatusInternalServerError(w, r, env)
		return
	}
	defer rows.Close()
	var eyecatches []EyeCatch
	for rows.Next() {
		var eyecatch EyeCatch
		err := rows.Scan(&eyecatch.Name, &eyecatch.Src)
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
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
