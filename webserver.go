package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Page struct {
	Title string
	Text  string
}

func index_handler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "This is a templated site", Text: "This is pretty neat"}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, p)
}

func about_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Written by Quinn, according a tutorial by Sentdex on Youtube")
}

func main() {
	http.HandleFunc("/", index_handler)
	http.HandleFunc("/about", about_handler)
	http.ListenAndServe(":8000", nil)
}
