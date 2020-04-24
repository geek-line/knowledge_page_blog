package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

const lenPathDelete = len("/admin/delete/")

func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

//AdminDeleteHandler admin/deleteに対するハンドラ
func AdminDeleteHandler(w http.ResponseWriter, r *http.Request, env map[string]string, db *sql.DB) {
	store := sessions.NewCookieStore([]byte(env["SESSION_KEY"]))
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
		return
	}
	suffix := r.URL.Path[lenPathDelete:]
	defer db.Close()
	var id int
	id, _ = strconv.Atoi(suffix)
	rows, err := db.Query("DELETE FROM knowledges WHERE id = ?", id)
	if err != nil {
		log.Print(err.Error())
		StatusInternalServerError(w, r, env)
		return
	}
	defer rows.Close()
	rows, err = db.Query("DELETE FROM knowledges_tags WHERE knowledge_id = ?", id)
	if err != nil {
		log.Print(err.Error())
		StatusInternalServerError(w, r, env)
		return
	}
	defer rows.Close()
	http.Redirect(w, r, "/admin/knowledges/", http.StatusFound)
}
