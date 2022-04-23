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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "This is a templated site", Text: "This is pretty neat"}
	t, _ := template.ParseFiles("tmpl/index.html")
	t.Execute(w, p)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/input.html")
	t.Execute(w, nil)
}

func returnHandler(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")

	p := &Ret{Content: string(content)}
	t, _ := template.ParseFiles("tmpl/output.html")
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/get/", getHandler)
	http.HandleFunc("/return/", returnHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8000", nil)
}
