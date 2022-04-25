package main

import (
	"html/template"
	"net/http"
)

type Page struct {
	Title string
	Text  string
}

type Ret struct {
	Content string
}

type Server struct {
	Name   string
	Owner  string
	Status string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/index.html")
	t.Execute(w, nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/create.html")
	t.Execute(w, nil)
}

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/console.html")
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/console/", consoleHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8000", nil)
}
