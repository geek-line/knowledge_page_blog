package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
)

// KnowledgeLikeHandler /knowledge/likeに対するハンドラ
func KnowledgeLikeHandler(w http.ResponseWriter, r *http.Request, env map[string]string) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		panic(err.Error())
	}
	db, err := sql.Open("mysql", env["SQL_ENV"])
	if err != nil {
		panic(err.Error())
	}
	if _, err := db.Query("UPDATE knowledges SET likes = likes + 1 WHERE id = ?", id); err != nil {
		panic(err.Error())
	}
}
