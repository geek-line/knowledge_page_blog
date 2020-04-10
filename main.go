package main

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
)

func main() {
	dir, _ := os.Getwd()
	fmt.Println(reflect.TypeOf(dir))
	http.Handle("/new/", http.StripPrefix("/new/", http.FileServer(http.Dir(dir+"/static"))))
	http.ListenAndServe(":3000", nil)
}
