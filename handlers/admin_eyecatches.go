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

//AdminEyecatchesHandler /admin/eyecatchesに対するハンドラ
func AdminEyecatchesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	header := newHeader(true)
	switch {
	case r.Method == "GET":
		eyecatches, err := models.GetAllEyecatches()
		if err != nil {
			log.Print(err.Error())
		}
		t := template.Must(template.ParseFiles("template/admin_eyecatches.html", "template/_header.html"))
		if err := t.Execute(w, struct {
			Header     structs.Header
			Eyecatches []structs.Eyecatch
		}{
			Header:     header,
			Eyecatches: eyecatches,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case r.Method == "POST":
		name := r.FormValue("name")
		src := r.FormValue("src")
		err := models.PostEyecatch(name, src)
		if err != nil {
			log.Print(err.Error())
			return
		}
		http.Redirect(w, r, routes.AdminEyecatchesPath, http.StatusFound)
	case r.Method == "PUT":
		id, _ := strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("name")
		src := r.FormValue("src")
		err := models.UpdateEyecatch(id, name, src)
		if err != nil {
			log.Print(err.Error())
			return
		}
	case r.Method == "DELETE":
		id, _ := strconv.Atoi(r.FormValue("id"))
		err := models.DeleteEyecatch(id)
		if err != nil {
			log.Print(err.Error())
			return
		}
	}
}
