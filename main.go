package main

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

type IndexPage struct {
	Id               int    //タイトル
	Title            string //タイトルの中身
	SelectedTagNames []string
	UpdatedAt        string
	Likes            int
}

type DetailPage struct {
	Id               int
	Title            string
	Content          string
	SelectedTagNames []string
	UpdatedAt        string
	Likes            int
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

	db, err := sql.Open("mysql", env["sqlEnv"])
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

func adminKnowledgesHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
	}
	header := newHeader(true)
	suffix := r.URL.Path[lenPathAdminKnowledges:]
	db, err := sql.Open("mysql", env["sqlEnv"])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if suffix != "" {
		var editPage Knowledges
		knowledgeID, _ := strconv.Atoi(suffix)
		err := db.QueryRow("SELECT id, title, content FROM knowledges WHERE id = ?", knowledgeID).Scan(&editPage.Id, &editPage.Title, &editPage.Content)
		switch {
		case err == sql.ErrNoRows:
			log.Println("レコードが存在しません")
			statusNotFoundHandler(w, r)
		default:
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
			var selectedTagsID []int

			rows, _ = db.Query("SELECT tag_id FROM knowledges_tags WHERE knowledge_id = ?", knowledgeID)
			for rows.Next() {
				var selectedTagID int
				err := rows.Scan(&selectedTagID)
				if err != nil {
					panic(err.Error())
				}
				selectedTagsID = append(selectedTagsID, selectedTagID)
			}

			t := template.Must(template.ParseFiles("template/admin_edit.html", "template/_header.html"))
			if err := t.Execute(w, struct {
				Header         Header
				EditPage       Knowledges
				Tags           []Tag
				SelectedTagsID []int
			}{
				Header:         header,
				EditPage:       editPage,
				Tags:           tags,
				SelectedTagsID: selectedTagsID,
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
		createdAt := time.Now()
		updatedAt := time.Now()
		stmtInsert, err := db.Prepare("INSERT INTO knowledges(title, content, created_at, updated_at) VALUES(?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		defer stmtInsert.Close()
		result, err := stmtInsert.Exec(title, content, createdAt, updatedAt)
		if err != nil {
			panic(err.Error())
		}
		knowledgeID, err := result.LastInsertId()
		if err != nil {
			panic(err.Error())
		}
		tags := strings.Split(r.FormValue("tags"), ",")
		for _, tag := range tags {
			tagID, _ := strconv.Atoi(tag)
			_, err = db.Query("INSERT INTO knowledges_tags(knowledge_id, tag_id, created_at, updated_at) VALUES(?, ?, ?, ?)", knowledgeID, tagID, createdAt, updatedAt)
			if err != nil {
				panic(err.Error())
			}
		}
	case r.Method == "PUT":
		knowledgeID, _ := strconv.Atoi(r.FormValue("id"))
		updatedAt := time.Now()
		_, err = db.Query("UPDATE knowledges SET title = ?, content = ?, updated_at = ? WHERE id = ?", title, content, updatedAt, knowledgeID)
		if err != nil {
			panic(err.Error())
		}
		_, err := db.Query("DELETE FROM knowledges_tags WHERE knowledge_id = ?", knowledgeID)
		if err != nil {
			panic(err.Error())
		}
		tags := strings.Split(r.FormValue("tags"), ",")
		createdAt := time.Now()
		for _, tag := range tags {
			tagID, _ := strconv.Atoi(tag)
			_, err = db.Query("INSERT INTO knowledges_tags(knowledge_id, tag_id, created_at, updated_at) VALUES(?, ?, ?, ?)", knowledgeID, tagID, createdAt, updatedAt)
			if err != nil {
				panic(err.Error())
			}
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
	_, err = db.Query("DELETE FROM knowledges_tags WHERE knowledge_id = ?", id)
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
		var detailPage DetailPage
		var id int
		id, _ = strconv.Atoi(suffix)
		err := db.QueryRow("SELECT id, title, content, updated_at, likes FROM knowledges WHERE id = ?", id).Scan(&detailPage.Id, &detailPage.Title, &detailPage.Content, &detailPage.UpdatedAt, &detailPage.Likes)
		switch {
		case err == sql.ErrNoRows:
			log.Println("レコードが存在しません")
			statusNotFoundHandler(w, r)
		case err != nil:
			panic(err.Error())
		default:
			var selectedTagNames []string
			tagsRows, err := db.Query("SELECT tag_id FROM knowledges_tags WHERE knowledge_id = ?", detailPage.Id)
			if err != nil {
				panic(err.Error())
			}
			defer tagsRows.Close()
			for tagsRows.Next() {
				var selectedTagID int
				var selectedTagName string
				err := tagsRows.Scan(&selectedTagID)
				if err != nil {
					panic(err.Error())
				}
				if err = db.QueryRow("SELECT name FROM tags WHERE id = ?", selectedTagID).Scan(&selectedTagName); err != nil {
					panic(err.Error())
				}
				selectedTagNames = append(selectedTagNames, selectedTagName)
			}
			detailPage.SelectedTagNames = selectedTagNames
			t := template.Must(template.ParseFiles("template/user_details.html", "template/_header.html", "template/_footer.html"))
			if err := t.Execute(w, struct {
				Header     Header
				DetailPage DetailPage
			}{
				Header:     header,
				DetailPage: detailPage,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		rows, err := db.Query("SELECT id, title, updated_at, likes FROM knowledges")
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		var indexPages []IndexPage
		for rows.Next() {
			var indexPage IndexPage
			err := rows.Scan(&indexPage.Id, &indexPage.Title, &indexPage.UpdatedAt, &indexPage.Likes)
			if err != nil {
				panic(err.Error())
			}
			var selectedTagNames []string
			tagsRows, err := db.Query("SELECT tag_id FROM knowledges_tags WHERE knowledge_id = ?", indexPage.Id)
			if err != nil {
				panic(err.Error())
			}
			defer tagsRows.Close()
			for tagsRows.Next() {
				var selectedTagID int
				var selectedTagName string
				err := tagsRows.Scan(&selectedTagID)
				if err != nil {
					panic(err.Error())
				}
				err = db.QueryRow("SELECT name FROM tags WHERE id = ?", selectedTagID).Scan(&selectedTagName)
				selectedTagNames = append(selectedTagNames, selectedTagName)
			}
			indexPage.SelectedTagNames = selectedTagNames
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

func statusNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	header := newHeader(false)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		header.IsLogin = true
	}
	t := template.Must(template.ParseFiles("template/404.html", "template/_header.html"))
	if err := t.Execute(w, struct {
		Header Header
	}{
		Header: header,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func knowledgeLikeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		panic(err.Error())
	}
	db, err := sql.Open("mysql", env["sqlEnv"])
	if err != nil {
		panic(err.Error())
	}
	if _, err := db.Query("UPDATE knowledges SET likes = likes + 1 WHERE id = ?", id); err != nil {
		panic(err.Error())
	}
}

func main() {
	dir, _ := os.Getwd()
	http.HandleFunc("/", statusNotFoundHandler)
	http.HandleFunc("/admin/login/", adminLoginHandler)
	http.HandleFunc("/admin/logout/", adminLogoutHandler)
	http.HandleFunc("/admin/knowledges/", adminKnowledgesHandler)
	http.HandleFunc("/admin/tags/", adminTagsHandler)
	http.HandleFunc("/admin/new/", adminNewHandler)
	http.HandleFunc("/admin/save/", adminSaveHandler)
	http.HandleFunc("/admin/delete/", adminDeleteHandler)
	http.HandleFunc("/knowledges/", knowledgesHandler)
	http.HandleFunc("/knowledges/like", knowledgeLikeHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir(dir+"/node_modules/"))))
	http.ListenAndServe(":3000", nil)
}
