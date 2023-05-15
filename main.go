package main

import (
	"net/http"
	"html/template"
    "log"
)

var tmpl = template.Must(template.ParseGlob("_includes/*.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "head.html",nil)
	tmpl.ExecuteTemplate(w, "index.html","Blogger")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/submit", uploadHandler)
	log.Println("Starting server at localhost:4000")
	if err := http.ListenAndServe("localhost:4000", mux); err != nil {
		log.Fatal(err)
	}
}

func chekIfAllFoldersExist() {
	
}