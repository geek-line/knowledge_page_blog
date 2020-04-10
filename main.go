package main

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Knowledges struct {
	ID      int
	Title   string
	Content string
}

func knowledgesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:Reibo1998@@/knowledge_blog")
	if err != nil {
		panic(err.Error())
	}
	rows, err := db.Query("select id title from knowledges")
	if err != nil {
		panic(err.Error())
	}
	var knowledges Knowledges
	var idsTitles [][]string

	for rows.Next() {
		err := rows.Scan(&knowledges.ID, &knowledges.Title)
		if err != nil {
			panic(err.Error())
		}
		idTitle := []string{string(knowledges.ID), knowledges.Title}
		idsTitles = append(idsTitles, idTitle)
	}
}

func main() {
	dir, _ := os.Getwd()
	http.HandleFunc("admin/knowledges/", knowledgesHandler)
	http.HandleFunc("admin/new")
	http.Handle("/new/", http.StripPrefix("/new/", http.FileServer(http.Dir(dir+"/static/"))))
	http.ListenAndServe(":3000", nil)
}
