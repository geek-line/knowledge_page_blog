package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"

	"../models"
	"../routes"
	"../structs"
)

const lenPathKnowledges = len(routes.UserKnowledgesPath)

//KnowledgesHandler /knowledgesに対するハンドラ
func KnowledgesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, auth bool) {
	header := newHeader(false)
	if auth {
		header.IsLogin = true
	}
	suffix := r.URL.Path[lenPathKnowledges:]
	if suffix == "" || suffix == "search" {
		pageNum := 1
		query := r.URL.Query()
		if query["page"] != nil {
			var err error
			if pageNum, err = strconv.Atoi(query.Get("page")); err != nil {
				StatusNotFoundHandler(w, r, auth)
				return
			}
		}
		sortKey := "updated_at"
		var currentSort string
		if query["sort"] != nil {
			switch {
			case query.Get("sort") == "create":
				sortKey = "created_at"
				currentSort = "create"
				break
			case query.Get("sort") == "update":
				sortKey = "updated_at"
				currentSort = "update"
				break
			case query.Get("sort") == "like":
				sortKey = "likes"
				currentSort = "like"
				break
			default:
				StatusNotFoundHandler(w, r, auth)
				break
			}
		} else {
			currentSort = "create"
		}
		tagRankingElem, err := models.GetTop10ReferencedTags()
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, auth)
			return
		}
		NumOfKnowledges, err := models.GetNumOfKnowledges()
		if err != nil {
			log.Print(err.Error())
			StatusNotFoundHandler(w, r, auth)
			return
		}
		pageNums := int(math.Ceil(NumOfKnowledges / 20))
		if pageNums < pageNum {
			StatusNotFoundHandler(w, r, auth)
			return
		}
		var pageNationElems = make([]structs.Page, pageNums)
		for i := 0; i < pageNums; i++ {
			pageNationElems[i].PageNum = i + 1
			pageNationElems[i].IsSelect = false
		}
		pageNationElems[pageNum-1].IsSelect = true
		var pageNation structs.PageNation
		pageNation.PageElems = pageNationElems
		pageNation.PageNum = pageNum
		pageNation.NextPageNum = pageNum + 1
		pageNation.PrevPageNum = pageNum - 1
		indexElems, err := models.Get20SortedElems(sortKey, (pageNum-1)*20, 20)
		if err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, auth)
			return
		}
		indexPage := structs.UserIndexPage{
			PageNation:  pageNation,
			IndexElems:  indexElems,
			CurrentSort: currentSort,
			TagRanking:  tagRankingElem,
		}
		t := template.Must(template.ParseFiles("template/user_knowledges.html", "template/_header.html", "template/_footer.html"))
		if err = t.Execute(w, struct {
			Header    structs.Header
			IndexPage structs.UserIndexPage
		}{
			Header:    header,
			IndexPage: indexPage,
		}); err != nil {
			log.Print(err.Error())
			StatusInternalServerError(w, r, auth)
		}
	} else {
		var userDetailPage structs.UserDetailPage
		var id int
		var err error
		if id, err = strconv.Atoi(suffix); err != nil {
			log.Print(err.Error())
			StatusNotFoundHandler(w, r, auth)
			return
		}
		userDetailPage.Knowledge, err = models.GetKnowledge(id)
		switch {
		case err == sql.ErrNoRows:
			log.Println("レコードが存在しません")
			StatusNotFoundHandler(w, r, auth)
		case err != nil:
			log.Print(err.Error())
			StatusInternalServerError(w, r, auth)
			return
		default:
			userDetailPage.SelectedTags, err = models.GetTagFromKnowledgeID(userDetailPage.Knowledge.ID)
			if err != nil {
				log.Print(err.Error())
				StatusInternalServerError(w, r, auth)
				return
			}
			t := template.Must(template.ParseFiles("template/user_details.html", "template/_header.html", "template/_footer.html"))
			if err := t.Execute(w, struct {
				Header     structs.Header
				DetailPage structs.UserDetailPage
			}{
				Header:     header,
				DetailPage: userDetailPage,
			}); err != nil {
				log.Print(err.Error())
				StatusInternalServerError(w, r, auth)
			}
		}
	}
}
