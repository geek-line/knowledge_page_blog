package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const lenPathTags = len("/tags/")

//TagsHandler /tags/に対するハンドラ
func TagsHandler(w http.ResponseWriter, r *http.Request, env map[string]string) {
	session, _ := store.Get(r, "cookie-name")
	header := newHeader(false)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		header.IsLogin = true
	}

	suffix := r.URL.Path[lenPathTags:]
	db, err := sql.Open("mysql", env["SQL_ENV"])
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if suffix != "" {
		id, err := strconv.Atoi(suffix)
		if err != nil {
			StatusNotFoundHandler(w, r)
			return
		}
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
		var filteredTagName string
		db.QueryRow("SELECT name FROM tags WHERE id = ?", id).Scan(&filteredTagName)
		rows, err = db.Query("SELECT knowledges.id, title, knowledges.updated_at, likes, eyecatch_src FROM knowledges INNER JOIN knowledges_tags ON knowledges_tags.knowledge_id = knowledges.id WHERE tag_id = ?", id)
		if err != nil {
			// panic(err.Error())
			StatusNotFoundHandler(w, r)
			log.Println("クエリのエラー")
			return
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
			tagsRows, err := db.Query("SELECT tags.id, tags.name FROM tags INNER JOIN knowledges_tags ON knowledges_tags.tag_id = tags.id WHERE knowledge_id = ?", indexElem.ID)
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
				// db.QueryRow("SELECT name FROM tags WHERE id = ?", selectedTag.ID).Scan(&selectedTag.Name)
				selectedTags = append(selectedTags, selectedTag)
			}
			indexElem.SelectedTags = selectedTags
			indexPage = append(indexPage, indexElem)
		}
		t := template.Must(template.ParseFiles("template/user_tags.html", "template/_header.html", "template/_footer.html"))
		if err = t.Execute(w, struct {
			Header          Header
			Tags            []Tag
			IndexPage       []IndexElem
			FilteredTagName string
		}{
			Header:          header,
			Tags:            tags,
			IndexPage:       indexPage,
			FilteredTagName: filteredTagName,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		StatusNotFoundHandler(w, r)
	}
}
