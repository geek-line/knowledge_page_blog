package handlers

import (
	"html/template"
	"log"
	"net/http"

	"../models"
	"../structs"
)

func newHeader(isLogin bool) structs.Header {
	return structs.Header{IsLogin: isLogin}
}

//AdminNewHandler /admin/newに対するハンドラ
func AdminNewHandler(w http.ResponseWriter, r *http.Request) {
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
	t := template.Must(template.ParseFiles("template/admin_new.html", "template/_header.html"))
	header := newHeader(true)
	if err := t.Execute(w, struct {
		Header     structs.Header
		Tags       []structs.Tag
		Eyecatches []structs.Eyecatch
	}{
		Header:     header,
		Tags:       tags,
		Eyecatches: eyecatches,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
