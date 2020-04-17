package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//AdminSaveHandler /admin/saveに対するハンドラ
func AdminSaveHandler(w http.ResponseWriter, r *http.Request, env map[string]string) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	eyecatchSrc := r.FormValue("eyecatch_src")
	db, err := sql.Open("mysql", env["SQL_ENV"])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	switch {
	case r.Method == "POST":
		createdAt := time.Now()
		updatedAt := time.Now()
		stmtInsert, err := db.Prepare("INSERT INTO knowledges(title, content, created_at, updated_at, eyecatch_src) VALUES(?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		defer stmtInsert.Close()
		result, err := stmtInsert.Exec(title, content, createdAt, updatedAt, eyecatchSrc)
		if err != nil {
			panic(err.Error())
		}
		knowledgeID, err := result.LastInsertId()
		if err != nil {
			panic(err.Error())
		}
		if r.FormValue("tags") != "" {
			tags := strings.Split(r.FormValue("tags"), ",")
			for _, tag := range tags {
				tagID, _ := strconv.Atoi(tag)
				_, err = db.Query("INSERT INTO knowledges_tags(knowledge_id, tag_id, created_at, updated_at) VALUES(?, ?, ?, ?)", knowledgeID, tagID, createdAt, updatedAt)
				if err != nil {
					panic(err.Error())
				}
			}
		}
	case r.Method == "PUT":
		knowledgeID, _ := strconv.Atoi(r.FormValue("id"))
		updatedAt := time.Now()
		_, err = db.Query("UPDATE knowledges SET title = ?, content = ?, updated_at = ?, eyecatch_src = ? WHERE id = ?", title, content, updatedAt, eyecatchSrc, knowledgeID)
		if err != nil {
			panic(err.Error())
		}
		_, err := db.Query("DELETE FROM knowledges_tags WHERE knowledge_id = ?", knowledgeID)
		if err != nil {
			panic(err.Error())
		}
		if r.FormValue("tags") != "" {
			tags := strings.Split(r.FormValue("tags"), ",")
			createdAt := time.Now()
			for _, tag := range tags {
				tagID, _ := strconv.Atoi(tag)
				_, err = db.Query("INSERT INTO knowledges_tags(knowledge_id, tag_id, created_at, updated_at) VALUES(?, ?, ?, ?)", knowledgeID, tagID, createdAt, updatedAt)
				if err != nil {
					panic(err.Error())
				}
			}
		}
		return
	default:
		break
	}
	http.Redirect(w, r, "/admin/knowledges/", http.StatusFound)
}