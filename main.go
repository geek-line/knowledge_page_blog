package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"os"

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

var templates = make(map[string]*template.Template)

func knowledgesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:Reibo1998@@/knowledge_blog")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
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

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	db, err := sql.Open("mysql", "root:Reibo1998@@/knowledge_blog")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Query("INSERT INTO knowledges(title, content) VALUES(?, ?)", title, content)
	if err != nil {
		panic(err.Error())
	}
	// lastInsertID, err := result.LastInsertId()
	// if err != nil {
	// 	panic(err.Error())
	// }
	// log.Println(lastInsertID)
	http.Redirect(w, r, "/admin/knowledges/", http.StatusFound)
}

func main() {
	dir, _ := os.Getwd()
	http.HandleFunc("/admin/knowledges/", knowledgesHandler)
	http.HandleFunc("/admin/save/", saveHandler)
	http.Handle("/admin/new/", http.StripPrefix("/admin/new/", http.FileServer(http.Dir(dir+"/static/"))))
	http.ListenAndServe(":3000", nil)
}
