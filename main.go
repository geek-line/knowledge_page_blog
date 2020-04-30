package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

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
	http.Redirect(w, r, userKnowledgesPath, http.StatusFound)
}

func init() {
	envLoad()
	env["SESSION_KEY"] = os.Getenv("SESSION_KEY")
	env["SQL_ENV"] = os.Getenv("SQL_ENV")
}

func main() {
	dir, _ := os.Getwd()
	http.HandleFunc(rootPath, redirectHandler)
	http.HandleFunc(adminLoginPath, makeHandlerUsingEnv(handlers.AdminLoginHandler))
	http.HandleFunc(amdinLogoutPath, makeHandlerUsingEnv(handlers.AdminLogoutHandler))
	http.HandleFunc(adminKnowledgesPath, makeHandlerUsingEnv(handlers.AdminKnowledgesHandler))
	http.HandleFunc(adminTagsPath, makeHandlerUsingEnv(handlers.AdminTagsHandler))
	http.HandleFunc(adminNewPath, makeHandlerUsingEnv(handlers.AdminNewHandler))
	http.HandleFunc(adminSavePath, makeHandlerUsingEnv(handlers.AdminSaveHandler))
	http.HandleFunc(adminDeletePath, makeHandlerUsingEnv(handlers.AdminDeleteHandler))
	http.HandleFunc(adminEyecatchesPath, makeHandlerUsingEnv(handlers.AdminEyeCatchesHandler))
	http.HandleFunc(userKnowledgesPath, makeHandlerUsingEnv(handlers.KnowledgesHandler))
	http.HandleFunc(userKnowledgesLikePath, makeHandlerUsingEnv(handlers.KnowledgeLikeHandler))
	http.HandleFunc(userTagsPath, makeHandlerUsingEnv(handlers.TagsHandler))
	http.Handle(staticPath, http.StripPrefix(staticPath, http.FileServer(http.Dir(dir+staticPath))))
	http.Handle(nodeModulesPath, http.StripPrefix(nodeModulesPath, http.FileServer(http.Dir(dir+nodeModulesPath))))
	http.Handle(googleSitemapPath, http.StripPrefix(googleSitemapPath, http.FileServer(http.Dir(dir+googleSitemapPath))))
	http.ListenAndServe(":3000", nil)
}
