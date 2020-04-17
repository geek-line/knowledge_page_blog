package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const lenPathKnowledges = len("/knowledges/")

//KnowledgesHandler /knowledgesに対するハンドラ
func KnowledgesHandler(w http.ResponseWriter, r *http.Request, env map[string]string) {
	session, _ := store.Get(r, "cookie-name")
	header := newHeader(false)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		header.IsLogin = true
	}

	suffix := r.URL.Path[lenPathKnowledges:]
	db, err := sql.Open("mysql", env["SQL_ENV"])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if suffix != "" {
		var detailPage DetailPage
		var id int
		id, _ = strconv.Atoi(suffix)
		err := db.QueryRow("SELECT id, title, content, updated_at, likes, eyecatch_src FROM knowledges WHERE id = ?", id).Scan(&detailPage.Id, &detailPage.Title, &detailPage.Content, &detailPage.UpdatedAt, &detailPage.Likes, &detailPage.EyeCatchSrc)
		switch {
		case err == sql.ErrNoRows:
			log.Println("レコードが存在しません")
			StatusNotFoundHandler(w, r)
		case err != nil:
			panic(err.Error())
		default:
			var selectedTagNames []string
			tagsRows, err := db.Query("SELECT tag_id FROM knowledges_tags WHERE knowledge_id = ?", detailPage.Id)
			if err != nil {
				panic(err.Error())
			}
			defer tagsRows.Close()
			for tagsRows.Next() {
				var selectedTagID int
				var selectedTagName string
				err := tagsRows.Scan(&selectedTagID)
				if err != nil {
					panic(err.Error())
				}
				if err = db.QueryRow("SELECT name FROM tags WHERE id = ?", selectedTagID).Scan(&selectedTagName); err != nil {
					panic(err.Error())
				}
				selectedTagNames = append(selectedTagNames, selectedTagName)
			}
			detailPage.SelectedTagNames = selectedTagNames
			t := template.Must(template.ParseFiles("template/user_details.html", "template/_header.html", "template/_footer.html"))
			if err := t.Execute(w, struct {
				Header     Header
				DetailPage DetailPage
			}{
				Header:     header,
				DetailPage: detailPage,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		rows, err := db.Query("SELECT id, title, updated_at, likes, eyecatch_src FROM knowledges")
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		var indexPages []IndexPage
		for rows.Next() {
			var indexPage IndexPage
			err := rows.Scan(&indexPage.Id, &indexPage.Title, &indexPage.UpdatedAt, &indexPage.Likes, &indexPage.EyeCatchSrc)
			if err != nil {
				panic(err.Error())
			}
			var selectedTagNames []string
			tagsRows, err := db.Query("SELECT tag_id FROM knowledges_tags WHERE knowledge_id = ?", indexPage.Id)
			if err != nil {
				panic(err.Error())
			}
			defer tagsRows.Close()
			for tagsRows.Next() {
				var selectedTagID int
				var selectedTagName string
				err := tagsRows.Scan(&selectedTagID)
				if err != nil {
					panic(err.Error())
				}
				err = db.QueryRow("SELECT name FROM tags WHERE id = ?", selectedTagID).Scan(&selectedTagName)
				selectedTagNames = append(selectedTagNames, selectedTagName)
			}
			indexPage.SelectedTagNames = selectedTagNames
			indexPages = append(indexPages, indexPage)
		}
		t := template.Must(template.ParseFiles("template/user_knowledges.html", "template/_header.html", "template/_footer.html"))
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
