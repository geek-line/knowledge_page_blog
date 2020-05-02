package handlers

import (
	"log"
	"net/http"
	"strconv"

	"../models"
)

//SeatchHandler /search/に対するハンドラ
func SeatchHandler(w http.ResponseWriter, r *http.Request, auth bool) {
	header := newHeader(false)
	if auth {
		header.IsLogin = true
	}
	pageNum := 1
	var err error
	query := r.URL.Query()
	if query["page"] != nil {
		if pageNum, err = strconv.Atoi(query.Get("page")); err != nil {
			StatusNotFoundHandler(w, r, auth)
			return
		}
	}
	sortKey := "created_at"
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
}
