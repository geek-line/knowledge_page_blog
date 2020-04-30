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
			// var selectedTagsID []int
			selectedTagsID, err := models.GetTagIDsFromKnowledgeID(knowledgeID)

			// rows, _ := db.Query("SELECT tag_id FROM knowledges_tags WHERE knowledge_id = ?", knowledgeID)
			// for rows.Next() {
			// 	var selectedTagID int
			// 	err := rows.Scan(&selectedTagID)
			if err != nil {
				log.Print(err.Error())
				return
			}
			// 	selectedTagsID = append(selectedTagsID, selectedTagID)
			// }

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
		// var indexPage IndexPage
		// var err error
		knowledges, err := models.GetAllKnowledges()
		// // rows, err := db.Query("SELECT id, title, created_at, updated_at FROM knowledges")
		if err != nil {
			log.Print(err.Error())
			return
		}
		// defer rows.Close()

		// for rows.Next() {
		// 	var indexElem IndexElem
		// 	err := rows.Scan(&indexElem.Knowledge.ID, &indexElem.Knowledge.Title, &indexElem.Knowledge.CreatedAt, &indexElem.Knowledge.UpdatedAt)
		// 	if err != nil {
		// 		log.Print(err.Error())
		// 		return
		// 	}
		// 	indexPage = append(indexPage, indexElem)
		// }

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
