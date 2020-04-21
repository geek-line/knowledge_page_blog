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

func makeHandlerUsingMySQL(fn func(w http.ResponseWriter, r *http.Request, env map[string]string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, env)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/knowledges/", http.StatusFound)
}

func init() {
	envLoad()
	env["SQL_ENV"] = os.Getenv("SQL_ENV")
}

func main() {
	dir, _ := os.Getwd()
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/admin/login/", makeHandlerUsingMySQL(handlers.AdminLoginHandler))
	http.HandleFunc("/admin/logout/", handlers.AdminLogoutHandler)
	http.HandleFunc("/admin/knowledges/", makeHandlerUsingMySQL(handlers.AdminKnowledgesHandler))
	http.HandleFunc("/admin/tags/", makeHandlerUsingMySQL(handlers.AdminTagsHandler))
	http.HandleFunc("/admin/new/", makeHandlerUsingMySQL(handlers.AdminNewHandler))
	http.HandleFunc("/admin/save/", makeHandlerUsingMySQL(handlers.AdminSaveHandler))
	http.HandleFunc("/admin/delete/", makeHandlerUsingMySQL(handlers.AdminDeleteHandler))
	http.HandleFunc("/knowledges/", makeHandlerUsingMySQL(handlers.KnowledgesHandler))
	http.HandleFunc("/knowledges/like", makeHandlerUsingMySQL(handlers.KnowledgeLikeHandler))
	http.HandleFunc("/tags/", makeHandlerUsingMySQL(handlers.TagsHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir(dir+"/node_modules/"))))
	http.ListenAndServe(":3000", nil)
	// l, err := net.Listen("tcp", "127.0.0.1:9000")
	// if err != nil {
	//     return
	// }
	// fcgi.Serve(l, nil)
}
