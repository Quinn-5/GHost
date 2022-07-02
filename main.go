package main

import (
	"html/template"
	"net/http"

	"github.com/Quinn-5/GHost/ghost"
	"github.com/Quinn-5/GHost/ghost/servconf"
	"k8s.io/apimachinery/pkg/api/resource"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/base.html", "tmpl/index.html")
	t.Execute(w, nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/base.html", "tmpl/create.html")

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
	if n, err := resource.ParseQuantity(r.FormValue("ram") + "Gi"); err == nil {
		ram = n
	}
	var disk resource.Quantity
	if n, err := resource.ParseQuantity(r.FormValue("disk") + "Gi"); err == nil {
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

	http.Redirect(w, r, "/result/", http.StatusFound)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/base.html", "tmpl/result.html")
	t.Execute(w, nil)
}

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/base.html", "tmpl/console.html")
	t.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/base.html", "tmpl/login.html")
	t.Execute(w, nil)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/create/", createHandler)
	mux.HandleFunc("/console/", consoleHandler)
	mux.HandleFunc("/login/", loginHandler)
	mux.HandleFunc("/result/", resultHandler)

	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.Handle("/tmpl/", http.StripPrefix("/tmpl/", http.FileServer(http.Dir("tmpl"))))

	http.ListenAndServe(":8000", mux)
}
