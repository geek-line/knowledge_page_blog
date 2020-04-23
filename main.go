package main

import (
	"log"
	"net/http"
	"os"

	"./handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var env = make(map[string]string)

func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func makeHandlerUsingEnv(fn func(w http.ResponseWriter, r *http.Request, env map[string]string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, env)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/knowledges/", http.StatusFound)
}

func init() {
	envLoad()
	env["SESSION_KEY"] = os.Getenv("SESSION_KEY")
	env["SQL_ENV"] = os.Getenv("SQL_ENV")
}

func main() {
	dir, _ := os.Getwd()
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/admin/login/", makeHandlerUsingEnv(handlers.AdminLoginHandler))
	http.HandleFunc("/admin/logout/", makeHandlerUsingEnv(handlers.AdminLogoutHandler))
	http.HandleFunc("/admin/knowledges/", makeHandlerUsingEnv(handlers.AdminKnowledgesHandler))
	http.HandleFunc("/admin/tags/", makeHandlerUsingEnv(handlers.AdminTagsHandler))
	http.HandleFunc("/admin/new/", makeHandlerUsingEnv(handlers.AdminNewHandler))
	http.HandleFunc("/admin/save/", makeHandlerUsingEnv(handlers.AdminSaveHandler))
	http.HandleFunc("/admin/delete/", makeHandlerUsingEnv(handlers.AdminDeleteHandler))
	http.HandleFunc("/knowledges/", makeHandlerUsingEnv(handlers.KnowledgesHandler))
	http.HandleFunc("/knowledges/like", makeHandlerUsingEnv(handlers.KnowledgeLikeHandler))
	http.HandleFunc("/tags/", makeHandlerUsingEnv(handlers.TagsHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir(dir+"/node_modules/"))))
	http.Handle("/google_sitemap/", http.StripPrefix("/google_sitemap/", http.FileServer(http.Dir(dir+"/google_sitemap/"))))
	http.ListenAndServe(":3000", nil)
}
