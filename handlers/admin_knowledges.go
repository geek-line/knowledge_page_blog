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

const lenPathAdminKnowledges = len(routes.AdminKnowledgesPath)

//AdminKnowledgesHandler admin/knowledgesに対するハンドラ
func AdminKnowledgesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	header := newHeader(true)
	suffix := r.URL.Path[lenPathAdminKnowledges:]
	if suffix != "" {
		var editPage structs.Knowledge
		knowledgeID, _ := strconv.Atoi(suffix)
		editPage, err := models.GetKnowledge(knowledgeID)
		switch {
		case err == sql.ErrNoRows:
			log.Print(err.Error())
		default:
			tags, err := models.GetAllTags()
			if err != nil {
				log.Print(err.Error())
				return
			}
			eyecatches, err := models.GetAllEyecatches()
			if err != nil {
				log.Print(err.Error())
				return
			}
			selectedTagsID, err := models.GetTagIDsFromKnowledgeID(knowledgeID)
			if err != nil {
				log.Print(err.Error())
				return
			}
			t := template.Must(template.ParseFiles("template/admin_edit.html", "template/_header.html"))
			if err := t.Execute(w, struct {
				Header         structs.Header
				EditPage       structs.Knowledge
				Tags           []structs.Tag
				Eyecatches     []structs.Eyecatch
				SelectedTagsID []int
			}{
				Header:         header,
				EditPage:       editPage,
				Tags:           tags,
				Eyecatches:     eyecatches,
				SelectedTagsID: selectedTagsID,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		knowledges, err := models.GetAllKnowledges()
		if err != nil {
			log.Print(err.Error())
			return
		}
		t := template.Must(template.ParseFiles("template/admin_knowledges.html", "template/_header.html"))
		header := newHeader(true)
		if err = t.Execute(w, struct {
			Header    structs.Header
			IndexPage []structs.Knowledge
		}{
			Header:    header,
			IndexPage: knowledges,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
