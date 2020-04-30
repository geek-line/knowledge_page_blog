package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

// KnowledgeLikeHandler /knowledge/likeに対するハンドラ
func KnowledgeLikeHandler(w http.ResponseWriter, r *http.Request, env map[string]string, db *sql.DB) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Print(err.Error())
		StatusInternalServerError(w, r, env)
		return
	}
	if r.Method == "POST" {
		if rows, err := db.Query("UPDATE knowledges SET likes = likes + 1 WHERE id = ?", id); err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
		} else {
			defer rows.Close()
		}
	} else {
		if rows, err := db.Query("UPDATE knowledges SET likes = likes - 1 WHERE id = ?", id); err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
		} else {
			defer rows.Close()
		}
	}
	return
}
