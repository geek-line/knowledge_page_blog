package main

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

type IndexPage struct {
	Id    int    //タイトル
	Title string //タイトルの中身
}

type Knowledges struct {
	Id      int
	Title   string
	Content string
}

type Header struct {
	IsLogin bool
}

type Tag struct {
	Id   int
	Name string
}

const (
	lenPathAdminKnowledges = len("/admin/knowledges/")
	lenPathDelete          = len("/admin/delete/")
	lenPathKnowledges      = len("/knowledges/")
)

var (
	env   = make(map[string]string)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func init() {

	sqlenv, err := ioutil.ReadFile("sql_env.txt")
	if err != nil {
		panic(err.Error())
	}
	env["sqlEnv"] = string(sqlenv)
}

func newHeader(isLogin bool) Header {
	return Header{IsLogin: isLogin}
}

func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
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
			http.Redirect(w, r, "/admin/knowledges/", http.StatusFound)
		}
	} else {
		email := r.FormValue("email")
		password := r.FormValue("password")
		db, err := sql.Open("mysql", env["sqlEnv"])
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		var correctPassword string
		if err = db.QueryRow("SELECT password FROM admin_user WHERE email = ?", email).Scan(&correctPassword); err != nil {
			http.Redirect(w, r, "/admin/login/", http.StatusFound)
		}
		if correctPassword == password {
			session, _ := store.Get(r, "cookie-name")
			session.Values["authenticated"] = true
			session.Save(r, w)
			http.Redirect(w, r, "/admin/knowledges/", http.StatusFound)

		} else {
			http.Redirect(w, r, "/admin/login/", http.StatusFound)
			return
		}
	}
}

func adminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/admin/login/", http.StatusFound)
}

func adminNewHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
	}
	t := template.Must(template.ParseFiles("template/admin_new.html", "template/_header.html"))
	header := newHeader(true)
	if err := t.Execute(w, struct {
		Header Header
	}{
		Header: header,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func adminKnowledgesHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
	}

	suffix := r.URL.Path[lenPathAdminKnowledges:]
	db, err := sql.Open("mysql", env["sqlEnv"])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if suffix != "" {
		var editPage Knowledges
		var id int
		id, _ = strconv.Atoi(suffix)
		err := db.QueryRow("SELECT id, title, content FROM knowledges WHERE id = ?", id).Scan(&editPage.Id, &editPage.Title, &editPage.Content)
		switch {
		case err == sql.ErrNoRows:
			log.Println("レコードが存在しません")
			http.NotFound(w, r)
		case err != nil:
			panic(err.Error())
		default:
			t := template.Must(template.ParseFiles("template/admin_edit.html", "template/_header.html"))
			header := newHeader(true)
			if err := t.Execute(w, struct {
				Header   Header
				EditPage Knowledges
			}{
				Header:   header,
				EditPage: editPage,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		rows, err := db.Query("SELECT id, title FROM knowledges")
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		var indexPages []IndexPage

		for rows.Next() {
			var indexPage IndexPage
			err := rows.Scan(&indexPage.Id, &indexPage.Title)
			if err != nil {
				panic(err.Error())
			}
			indexPages = append(indexPages, indexPage)
		}

		t := template.Must(template.ParseFiles("template/admin_knowledges.html", "template/_header.html"))
		header := newHeader(true)
		if err = t.Execute(w, struct {
			Header     Header
			IndexPages []IndexPage
		}{
			Header:     header,
			IndexPages: indexPages,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func adminSaveHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	db, err := sql.Open("mysql", env["sqlEnv"])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	switch {
	case r.Method == "POST":
		_, err = db.Query("INSERT INTO knowledges(title, content) VALUES(?, ?)", title, content)
		if err != nil {
			panic(err.Error())
		}
	case r.Method == "PUT":
		id := r.FormValue("id")
		_, err = db.Query("UPDATE knowledges SET title = ?, content = ? WHERE id = ?", title, content, id)
		if err != nil {
			panic(err.Error())
		}
		return
	default:
		break
	}
	http.Redirect(w, r, "/admin/knowledges/", http.StatusFound)
}

func adminDeleteHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
	}

	suffix := r.URL.Path[lenPathDelete:]
	db, err := sql.Open("mysql", env["sqlEnv"])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var id int
	id, _ = strconv.Atoi(suffix)
	_, err = db.Query("DELETE FROM knowledges WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}
	http.Redirect(w, r, "/admin/knowledges/", http.StatusFound)
}

func adminTagsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	header := newHeader(false)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		header.IsLogin = true
	}
	db, err := sql.Open("mysql", env["sqlEnv"])
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
			err := rows.Scan(&tag.Id, &tag.Name)
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
		var id int
		id, _ = strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("name")
		updatedAt := time.Now()
		_, err = db.Query("UPDATE tags SET name = ?, updated_at = ? WHERE id = ?", name, updatedAt, id)
		if err != nil {
			panic(err.Error())
		}
	}
}

func knowledgesHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	header := newHeader(false)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		header.IsLogin = true
	}

	suffix := r.URL.Path[lenPathKnowledges:]
	db, err := sql.Open("mysql", env["sqlEnv"])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if suffix != "" {
		var editPage Knowledges
		var id int
		id, _ = strconv.Atoi(suffix)
		err := db.QueryRow("SELECT id, title, content FROM knowledges WHERE id = ?", id).Scan(&editPage.Id, &editPage.Title, &editPage.Content)
		switch {
		case err == sql.ErrNoRows:
			log.Println("レコードが存在しません")
			http.NotFound(w, r)
		case err != nil:
			panic(err.Error())
		default:
			t := template.Must(template.ParseFiles("template/user_details.html", "template/_header.html", "template/_footer.html"))
			if err := t.Execute(w, struct {
				Header   Header
				EditPage Knowledges
			}{
				Header:   header,
				EditPage: editPage,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		rows, err := db.Query("SELECT id, title FROM knowledges")
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		var indexPages []IndexPage

		for rows.Next() {
			var indexPage IndexPage
			err := rows.Scan(&indexPage.Id, &indexPage.Title)
			if err != nil {
				panic(err.Error())
			}
			indexPages = append(indexPages, indexPage)
		}

		t := template.Must(template.ParseFiles("template/user_knowledges.html", "template/_header.html", "template/_footer.html"))
		if err = t.Execute(w, struct {
			Header     Header
			IndexPages []IndexPage
		}{
			Header:     header,
			IndexPages: indexPages,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}

func main() {
	dir, _ := os.Getwd()
	http.HandleFunc("/admin/login/", adminLoginHandler)
	http.HandleFunc("/admin/logout/", adminLogoutHandler)
	http.HandleFunc("/admin/knowledges/", adminKnowledgesHandler)
	http.HandleFunc("/admin/tags/", adminTagsHandler)
	http.HandleFunc("/admin/new/", adminNewHandler)
	http.HandleFunc("/admin/save/", adminSaveHandler)
	http.HandleFunc("/admin/delete/", adminDeleteHandler)
	http.HandleFunc("/knowledges/", knowledgesHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir(dir+"/node_modules/"))))
	http.ListenAndServe(":3000", nil)
}
