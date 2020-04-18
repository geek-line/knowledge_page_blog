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
		err := db.QueryRow("SELECT id, title, content, updated_at, likes, eyecatch_src FROM knowledges WHERE id = ?", id).Scan(&detailPage.ID, &detailPage.Title, &detailPage.Content, &detailPage.UpdatedAt, &detailPage.Likes, &detailPage.EyeCatchSrc)
		switch {
		case err == sql.ErrNoRows:
			log.Println("レコードが存在しません")
			StatusNotFoundHandler(w, r)
		case err != nil:
			panic(err.Error())
		default:
			var selectedTags []Tag
			tagsRows, err := db.Query("SELECT tags.id, tags.name FROM tags INNER JOIN knowledges_tags ON knowledges_tags.tag_id = tags.id WHERE knowledge_id = ?", detailPage.ID)
			if err != nil {
				panic(err.Error())
			}
			defer tagsRows.Close()
			for tagsRows.Next() {
				var selectedTag Tag
				err := tagsRows.Scan(&selectedTag.ID, &selectedTag.Name)
				if err != nil {
					panic(err.Error())
				}
				selectedTags = append(selectedTags, selectedTag)
			}
			detailPage.SelectedTags = selectedTags
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
		rows, err := db.Query("SELECT id, name FROM tags")
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		var tags []Tag
		for rows.Next() {
			var tag Tag
			err := rows.Scan(&tag.ID, &tag.Name)
			if err != nil {
				panic(err.Error())
			}
			tags = append(tags, tag)
		}
		rows, err = db.Query("SELECT id, title, updated_at, likes, eyecatch_src FROM knowledges")
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		var indexPage []IndexElem
		for rows.Next() {
			var indexElem IndexElem
			err := rows.Scan(&indexElem.ID, &indexElem.Title, &indexElem.UpdatedAt, &indexElem.Likes, &indexElem.EyeCatchSrc)
			if err != nil {
				panic(err.Error())
			}
			var selectedTags []Tag
			tagsRows, err := db.Query("SELECT tag_id FROM knowledges_tags WHERE knowledge_id = ?", indexElem.ID)
			if err != nil {
				panic(err.Error())
			}
			defer tagsRows.Close()
			for tagsRows.Next() {
				var selectedTag Tag
				err := tagsRows.Scan(&selectedTag.ID)
				if err != nil {
					panic(err.Error())
				}
				db.QueryRow("SELECT name FROM tags WHERE id = ?", selectedTag.ID).Scan(&selectedTag.Name)
				selectedTags = append(selectedTags, selectedTag)
			}
			indexElem.SelectedTags = selectedTags
			indexPage = append(indexPage, indexElem)
		}
		t := template.Must(template.ParseFiles("template/user_knowledges.html", "template/_header.html", "template/_footer.html"))
		if err = t.Execute(w, struct {
			Header    Header
			Tags      []Tag
			IndexPage []IndexElem
		}{
			Header:    header,
			Tags:      tags,
			IndexPage: indexPage,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
