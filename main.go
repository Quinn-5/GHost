package main

import (
	"html/template"
	"net/http"

	"github.com/Quinn-5/learning-go/ghost"
	"github.com/Quinn-5/learning-go/ghost/servconf"
	"k8s.io/apimachinery/pkg/api/resource"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/index.html")
	t.Execute(w, nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/create.html")

	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	servername := r.FormValue("servername")
	servertype := r.FormValue("servertype")
	var cpu resource.Quantity
	if n, err := resource.ParseQuantity(r.FormValue("cpu")); err == nil {
		cpu = n
	}
	var ram resource.Quantity
	if n, err := resource.ParseQuantity(r.FormValue("ram")); err == nil {
		ram = n
	}
	var disk resource.Quantity
	if n, err := resource.ParseQuantity(r.FormValue("disk")); err == nil {
		disk = n
	}

	p := &servconf.ServerConfig{
		Username:   username,
		Servername: servername,
		Type:       servertype,
		CPU:        cpu,
		RAM:        ram,
		Disk:       disk,
	}

	err := ghost.Create(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/console.html")
	t.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/login.html")
	t.Execute(w, nil)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/create/", createHandler)
	mux.HandleFunc("/console/", consoleHandler)
	mux.HandleFunc("/login/", loginHandler)

	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8000", mux)
}
