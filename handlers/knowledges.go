package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

const lenPathKnowledges = len("/knowledges/")

//KnowledgesHandler /knowledgesに対するハンドラ
func KnowledgesHandler(w http.ResponseWriter, r *http.Request, env map[string]string, db *sql.DB) {
	store := sessions.NewCookieStore([]byte(env["SESSION_KEY"]))
	session, _ := store.Get(r, "cookie-name")
	header := newHeader(false)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		header.IsLogin = true
	}
	suffix := r.URL.Path[lenPathKnowledges:]
	if suffix == "" || suffix == "search" {
		var indexPage IndexPage
		pageNum := 1
		query := r.URL.Query()
		if query["page"] != nil {
			var err error
			if pageNum, err = strconv.Atoi(query.Get("page")); err != nil {
				StatusNotFoundHandler(w, r, env)
				return
			}
		}
		sortKey := "updated_at"
		if query["sort"] != nil {
			switch {
			case query.Get("sort") == "update":
				sortKey = "updated_at"
				indexPage.CurrentSort = "update"
				break
			case query.Get("sort") == "like":
				sortKey = "likes"
				indexPage.CurrentSort = "like"
				break
			default:
				StatusNotFoundHandler(w, r, env)
				break
			}
		} else {
			indexPage.CurrentSort = "update"
		}
		rows, err := db.Query("SELECT tags.id, tags.name, count(knowledges_tags.id) AS count FROM tags INNER JOIN knowledges_tags ON knowledges_tags.tag_id = tags.id GROUP BY knowledges_tags.tag_id ORDER BY count DESC LIMIT 10;")
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		}
		defer rows.Close()
		var tags []Tag
		for rows.Next() {
			var tag Tag
			err := rows.Scan(&tag.ID, &tag.Name, &tag.CountOfUse)
			if err != nil {
				log.Print(err.Error())

				StatusInternalServerError(w, r, env)
				return
			}
			tags = append(tags, tag)
		}
		var knowledgeNums float64
		db.QueryRow("SELECT count(id) FROM knowledges").Scan(&knowledgeNums)
		pageNums := int(math.Ceil(knowledgeNums / 20))
		if pageNums < pageNum {
			StatusNotFoundHandler(w, r, env)
			return
		}
		var pageNationElems = make([]Page, pageNums)
		for i := 0; i < pageNums; i++ {
			pageNationElems[i].PageNum = i + 1
			pageNationElems[i].IsSelect = false
		}
		pageNationElems[pageNum-1].IsSelect = true
		indexPage.PageNation.PageElems = pageNationElems
		indexPage.PageNation.PageNum = pageNum
		indexPage.PageNation.NextPageNum = pageNum + 1
		indexPage.PageNation.PrevPageNum = pageNum - 1
		qtext := fmt.Sprintf("SELECT id, title, updated_at, likes, eyecatch_src FROM knowledges ORDER BY %s DESC LIMIT ?, ?", sortKey)
		rows, err = db.Query(qtext, (pageNum-1)*20, 20)
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var indexElem IndexElem
			err := rows.Scan(&indexElem.ID, &indexElem.Title, &indexElem.UpdatedAt, &indexElem.Likes, &indexElem.EyeCatchSrc)
			if err != nil {
				StatusNotFoundHandler(w, r, env)
				return
			}
			var selectedTags []Tag
			tagsRows, err := db.Query("SELECT tag_id FROM knowledges_tags WHERE knowledge_id = ?", indexElem.ID)
			if err != nil {
				log.Print(err.Error())
				StatusInternalServerError(w, r, env)
				return
			}
			defer tagsRows.Close()
			for tagsRows.Next() {
				var selectedTag Tag
				err := tagsRows.Scan(&selectedTag.ID)
				if err != nil {
					log.Print(err.Error())
					StatusInternalServerError(w, r, env)
					return
				}
				db.QueryRow("SELECT name FROM tags WHERE id = ?", selectedTag.ID).Scan(&selectedTag.Name)
				selectedTags = append(selectedTags, selectedTag)
			}
			indexElem.SelectedTags = selectedTags
			indexPage.IndexElems = append(indexPage.IndexElems, indexElem)
		}
		t := template.Must(template.ParseFiles("template/user_knowledges.html", "template/_header.html", "template/_footer.html"))
		if err = t.Execute(w, struct {
			Header    Header
			Tags      []Tag
			IndexPage IndexPage
		}{
			Header:    header,
			Tags:      tags,
			IndexPage: indexPage,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	} else {
		var detailPage DetailPage
		var id int
		var err error
		if id, err = strconv.Atoi(suffix); err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		}
		err = db.QueryRow("SELECT id, title, content, updated_at, likes, eyecatch_src FROM knowledges WHERE id = ?", id).Scan(&detailPage.ID, &detailPage.Title, &detailPage.Content, &detailPage.UpdatedAt, &detailPage.Likes, &detailPage.EyeCatchSrc)
		switch {
		case err == sql.ErrNoRows:
			log.Println("レコードが存在しません")
			StatusNotFoundHandler(w, r, env)
		case err != nil:
			log.Print(err.Error())
			StatusInternalServerError(w, r, env)
			return
		default:
			var selectedTags []Tag
			tagsRows, err := db.Query("SELECT tags.id, tags.name FROM tags INNER JOIN knowledges_tags ON knowledges_tags.tag_id = tags.id WHERE knowledge_id = ?", detailPage.ID)
			if err != nil {
				log.Print(err.Error())
				StatusInternalServerError(w, r, env)
				return
			}
			defer tagsRows.Close()
			for tagsRows.Next() {
				var selectedTag Tag
				err := tagsRows.Scan(&selectedTag.ID, &selectedTag.Name)
				if err != nil {
					log.Print(err.Error())
					StatusInternalServerError(w, r, env)
					return
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
	}
}
