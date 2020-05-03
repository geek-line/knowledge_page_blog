package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"../models"
	"../routes"
	"../structs"
)

const lenPathKnowledge = len(routes.UserKnowledgePath)

//KnowledgeHandler /knowledgesに対するハンドラ
func KnowledgeHandler(w http.ResponseWriter, r *http.Request, auth bool) {
	header := newHeader(false)
	if auth {
		header.IsLogin = true
	}
	suffix := r.URL.Path[lenPathKnowledge:]
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
