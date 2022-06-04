package main

import (
	"html/template"
	"net/http"

	"github.com/Quinn-5/learning-go/ghost/deployments"
	"k8s.io/apimachinery/pkg/api/resource"
)

// type serverType string

// const (
// 	CSGO      serverType = "csgo"
// 	Factorio             = "factorio"
// 	Minecraft            = "minecraft"
// 	TF2                  = "tf2"
// 	Terraria             = "terraria"
// )

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/index.html")
	t.Execute(w, nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	servername := r.FormValue("servername")
	servertype := r.FormValue("servertype")
	cpu := r.FormValue("cpu")
	ram := r.FormValue("ram")
	disk := r.FormValue("disk")

	p := &deployments.ServerConfig{
		Username:   username,
		Servername: servername,
		Type:       servertype,
		CPU:        resource.Format(cpu),
		RAM:        resource.Format(ram),
		Disk:       resource.Format(disk),
	}

	err := p.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, _ := template.ParseFiles("tmpl/create.html")
	t.Execute(w, nil)
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
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/console/", consoleHandler)
	http.HandleFunc("/login/", consoleHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8000", nil)
}
