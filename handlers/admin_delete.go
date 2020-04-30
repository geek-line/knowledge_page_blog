package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"../routes"
)

const lenPathDelete = len(routes.AdminDeletePath)

//AdminDeleteHandler admin/deleteに対するハンドラ
func AdminDeleteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	suffix := r.URL.Path[lenPathDelete:]
	defer db.Close()
	var id int
	id, _ = strconv.Atoi(suffix)
	rows, err := db.Query("DELETE FROM knowledges WHERE id = ?", id)
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer rows.Close()
	rows, err = db.Query("DELETE FROM knowledges_tags WHERE knowledge_id = ?", id)
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer rows.Close()
	http.Redirect(w, r, routes.AdminKnowledgesPath, http.StatusFound)
}
