package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"./routes"

	"./handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var env = make(map[string]string)
var db sql.DB

func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func makeHandlerUsingEnv(fn func(w http.ResponseWriter, r *http.Request, env map[string]string, db *sql.DB)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("mysql", env["SQL_ENV"])
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		fn(w, r, env, db)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, routes.UserKnowledgesPath, http.StatusFound)
}

func init() {
	envLoad()
	env["SESSION_KEY"] = os.Getenv("SESSION_KEY")
	env["SQL_ENV"] = os.Getenv("SQL_ENV")
}

func main() {
	dir, _ := os.Getwd()
	http.HandleFunc(routes.RootPath, redirectHandler)
	http.HandleFunc(routes.AdminLoginPath, makeHandlerUsingEnv(handlers.AdminLoginHandler))
	http.HandleFunc(routes.AdminLogoutPath, makeHandlerUsingEnv(handlers.AdminLogoutHandler))
	http.HandleFunc(routes.AdminKnowledgesPath, makeHandlerUsingEnv(handlers.AdminKnowledgesHandler))
	http.HandleFunc(routes.AdminTagsPath, makeHandlerUsingEnv(handlers.AdminTagsHandler))
	http.HandleFunc(routes.AdminNewPath, makeHandlerUsingEnv(handlers.AdminNewHandler))
	http.HandleFunc(routes.AdminSavePath, makeHandlerUsingEnv(handlers.AdminSaveHandler))
	http.HandleFunc(routes.AdminDeletePath, makeHandlerUsingEnv(handlers.AdminDeleteHandler))
	http.HandleFunc(routes.AdminEyecatchesPath, makeHandlerUsingEnv(handlers.AdminEyeCatchesHandler))
	http.HandleFunc(routes.UserKnowledgesPath, makeHandlerUsingEnv(handlers.KnowledgesHandler))
	http.HandleFunc(routes.UserKnowledgesLikePath, makeHandlerUsingEnv(handlers.KnowledgeLikeHandler))
	http.HandleFunc(routes.UserTagsPath, makeHandlerUsingEnv(handlers.TagsHandler))
	http.Handle(routes.StaticPath, http.StripPrefix(routes.StaticPath, http.FileServer(http.Dir(dir+routes.StaticPath))))
	http.Handle(routes.NodeModulesPath, http.StripPrefix(routes.NodeModulesPath, http.FileServer(http.Dir(dir+routes.NodeModulesPath))))
	http.Handle(routes.GoogleSitemapPath, http.StripPrefix(routes.GoogleSitemapPath, http.FileServer(http.Dir(dir+routes.GoogleSitemapPath))))
	http.ListenAndServe(":3000", nil)
}
