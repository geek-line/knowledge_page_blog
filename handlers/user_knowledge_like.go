package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"../models"
)

// KnowledgeLikeHandler /knowledge/likeに対するハンドラ
func KnowledgeLikeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, auth bool) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Print(err.Error())
		StatusInternalServerError(w, r, auth)
		return
	}
	if r.Method == "POST" {
		if err := models.IncrementLikes(id); err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, auth)
		}
	} else {
		if err := models.DecrementLikes(id); err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, auth)
		}
	}
	return
}
