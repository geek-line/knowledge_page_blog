package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"

	"../routes"
)

const lenPathTags = len(routes.UserTagsPath)

//TagsHandler /tags/に対するハンドラ
func TagsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, auth bool) {
	header := newHeader(false)
	if auth {
		header.IsLogin = true
	}
	suffix := r.URL.Path[lenPathTags:]
	if suffix != "" {
		var indexPage IndexPage
		pageNum := 1
		var err error
		query := r.URL.Query()
		if query["page"] != nil {
			if pageNum, err = strconv.Atoi(query.Get("page")); err != nil {
				StatusNotFoundHandler(w, r, auth)
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
				StatusNotFoundHandler(w, r, auth)
				break
			}
		} else {
			indexPage.CurrentSort = "update"
		}
		var filteredTag Tag
		filteredTag.ID, err = strconv.Atoi(suffix)
		if err != nil {
			StatusNotFoundHandler(w, r, auth)
			return
		}
		rows, err := db.Query("SELECT tags.id, tags.name, count(knowledges_tags.id) AS count FROM tags INNER JOIN knowledges_tags ON knowledges_tags.tag_id = tags.id GROUP BY knowledges_tags.tag_id ORDER BY count DESC LIMIT 10;")
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, auth)
			return
		}
		defer rows.Close()
		var tags []Tag
		for rows.Next() {
			var tag Tag
			err := rows.Scan(&tag.ID, &tag.Name, &tag.CountOfUse)
			if err != nil {
				log.Print(err.Error())
				StatusInternalServerError(w, r, auth)
				return
			}
			tags = append(tags, tag)
		}
		var knowledgeNums float64
		db.QueryRow("SELECT count(knowledges.id) FROM knowledges INNER JOIN knowledges_tags ON knowledges_tags.knowledge_id = knowledges.id WHERE tag_id = ?", filteredTag.ID).Scan(&knowledgeNums)
		pageNums := int(math.Ceil(knowledgeNums / 20))
		if pageNums < pageNum {
			StatusNotFoundHandler(w, r, auth)
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

		db.QueryRow("SELECT name FROM tags WHERE id = ?", filteredTag.ID).Scan(&filteredTag.Name)
		qtext := fmt.Sprintf("SELECT knowledges.id, title, knowledges.updated_at, likes, eyecatch_src FROM knowledges INNER JOIN knowledges_tags ON knowledges_tags.knowledge_id = knowledges.id WHERE tag_id = ? ORDER BY %s DESC LIMIT ?, ?", sortKey)
		rows, err = db.Query(qtext, filteredTag.ID, (pageNum-1)*20, 20)
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, auth)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var indexElem IndexElem
			err := rows.Scan(&indexElem.ID, &indexElem.Title, &indexElem.UpdatedAt, &indexElem.Likes, &indexElem.EyeCatchSrc)
			if err != nil {
				log.Print(err.Error())
				StatusInternalServerError(w, r, auth)
				return
			}
			var selectedTags []Tag
			tagsRows, err := db.Query("SELECT tags.id, tags.name FROM tags INNER JOIN knowledges_tags ON knowledges_tags.tag_id = tags.id WHERE knowledge_id = ?", indexElem.ID)
			if err != nil {
				log.Print(err.Error())
				StatusInternalServerError(w, r, auth)
				return
			}
			defer tagsRows.Close()
			for tagsRows.Next() {
				var selectedTag Tag
				err := tagsRows.Scan(&selectedTag.ID, &selectedTag.Name)
				if err != nil {
					log.Print(err.Error())
					StatusInternalServerError(w, r, auth)
					return
				}
				selectedTags = append(selectedTags, selectedTag)
			}
			indexElem.SelectedTags = selectedTags
			indexPage.IndexElems = append(indexPage.IndexElems, indexElem)
		}
		t := template.Must(template.ParseFiles("template/user_tags.html", "template/_header.html", "template/_footer.html"))
		if err = t.Execute(w, struct {
			Header      Header
			Tags        []Tag
			IndexPage   IndexPage
			FilteredTag Tag
		}{
			Header:      header,
			Tags:        tags,
			IndexPage:   indexPage,
			FilteredTag: filteredTag,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		StatusNotFoundHandler(w, r, auth)
	}
}
