package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const lenPathAdminKnowledges = len("/admin/knowledges/")

//AdminKnowledgesHandler admin/knowledgesに対するハンドラ
func AdminKnowledgesHandler(w http.ResponseWriter, r *http.Request, env map[string]string) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
	}
	header := newHeader(true)
	suffix := r.URL.Path[lenPathAdminKnowledges:]
	db, err := sql.Open("mysql", env["SQL_ENV"])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if suffix != "" {
		var editPage Knowledges
		knowledgeID, _ := strconv.Atoi(suffix)
		err := db.QueryRow("SELECT id, title, content, eyecatch_src FROM knowledges WHERE id = ?", knowledgeID).Scan(&editPage.Id, &editPage.Title, &editPage.Content, &editPage.EyeCatchSrc)
		switch {
		case err == sql.ErrNoRows:
			log.Println("レコードが存在しません")
			StatusNotFoundHandler(w, r)
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
