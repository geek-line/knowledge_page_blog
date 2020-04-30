package main

import (
	"net/http"
	"net/http/fcgi"
	"os"

	"./handlers"
	"./middleware"
	"./routes"
	_ "github.com/go-sql-driver/mysql"
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, routes.UserKnowledgesPath, http.StatusFound)
}

func main() {
	dir, _ := os.Getwd()
	http.HandleFunc(routes.RootPath, redirectHandler)
	http.HandleFunc(routes.AdminLoginPath, middleware.UserAuth(handlers.AdminLoginHandler))
	http.HandleFunc(routes.AdminLogoutPath, middleware.AdminAuth(handlers.AdminLogoutHandler))
	http.HandleFunc(routes.AdminKnowledgesPath, middleware.AdminAuth(handlers.AdminKnowledgesHandler))
	http.HandleFunc(routes.AdminTagsPath, middleware.AdminAuth(handlers.AdminTagsHandler))
	http.HandleFunc(routes.AdminNewPath, middleware.AdminAuth(handlers.AdminNewHandler))
	http.HandleFunc(routes.AdminSavePath, middleware.AdminAuth(handlers.AdminSaveHandler))
	http.HandleFunc(routes.AdminDeletePath, middleware.AdminAuth(handlers.AdminDeleteHandler))
	http.HandleFunc(routes.AdminEyecatchesPath, middleware.AdminAuth(handlers.AdminEyecatchesHandler))
	http.HandleFunc(routes.UserKnowledgesPath, middleware.UserAuth(handlers.KnowledgesHandler))
	http.HandleFunc(routes.UserKnowledgesLikePath, middleware.UserAuth(handlers.KnowledgeLikeHandler))
	http.HandleFunc(routes.UserTagsPath, middleware.UserAuth(handlers.TagsHandler))
	http.Handle(routes.StaticPath, http.StripPrefix(routes.StaticPath, http.FileServer(http.Dir(dir+routes.StaticPath))))
	http.Handle(routes.NodeModulesPath, http.StripPrefix(routes.NodeModulesPath, http.FileServer(http.Dir(dir+routes.NodeModulesPath))))
	http.Handle(routes.GoogleSitemapPath, http.StripPrefix(routes.GoogleSitemapPath, http.FileServer(http.Dir(dir+routes.GoogleSitemapPath))))
	http.ListenAndServe(":3000", nil)
	l, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		return
	}
	fcgi.Serve(l, nil)
}
