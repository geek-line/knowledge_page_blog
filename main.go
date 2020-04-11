package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
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

const lenPath = len("/admin/knowledges/")

func knowledgesHandler(w http.ResponseWriter, r *http.Request) {

	suffix := r.URL.Path[lenPath:]
	log.Println(suffix)
	db, err := sql.Open("mysql", "root:Reibo1998@@/knowledge_blog")
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
			t := template.Must(template.ParseFiles("template/admin_edit.html"))
			err = t.Execute(w, editPage)
			if err != nil {
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

		t := template.Must(template.ParseFiles("template/admin_knowledges.html"))
		err = t.Execute(w, indexPages)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	db, err := sql.Open("mysql", "root:Reibo1998@@/knowledge_blog")
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
	}
	http.Redirect(w, r, "/admin/knowledges/", http.StatusFound)
}

func main() {
	dir, _ := os.Getwd()
	http.HandleFunc("/admin/knowledges/", knowledgesHandler)
	http.HandleFunc("/admin/save/", saveHandler)
	http.Handle("/admin/new/", http.StripPrefix("/admin/new/", http.FileServer(http.Dir(dir+"/static/"))))
	http.ListenAndServe(":3000", nil)
}
