package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/Quinn-5/learning-go/ghost/deployments"
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
	// username := r.FormValue("username")
	// servername := r.FormValue("servername")
	// servertype := r.FormValue("servertype")
	// cpu := r.FormValue("cpu")
	// ram := r.FormValue("ram")
	// disk := r.FormValue("disk")

	// p := &deployments.ServerConfig{
	// 	Username:   username,
	// 	Servername: servername,
	// 	Type:       servertype,
	// 	CPU:        resource.Format(cpu),
	// 	RAM:        resource.Format(ram),
	// 	Disk:       resource.Format(disk),
	// }

	// err := p.Create()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

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
	// http.HandleFunc("/", indexHandler)
	// http.HandleFunc("/create/", createHandler)
	// http.HandleFunc("/console/", consoleHandler)
	// http.HandleFunc("/login/", loginHandler)

	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// http.Handle("/tmpl/", http.StripPrefix("/tmpl/", http.FileServer(http.Dir("tmpl"))))

	// http.ListenAndServe(":8000", nil)

	username := "Quinn"
	servername := "candice"
	servertype := "Minecraft"
	cpu := "1"
	var icpu int64
	if n, err := strconv.ParseInt(cpu, 10, 64); err == nil {
		icpu = n
	}
	ram := "4"
	var iram int64
	if n, err := strconv.ParseInt(ram, 10, 64); err == nil {
		iram = n
	}
	disk := "1024"
	var idisk int64
	if n, err := strconv.ParseInt(disk, 10, 64); err == nil {
		idisk = n
	}
	println(idisk)

	p := &deployments.ServerConfig{
		Username:   strings.ToLower(username),
		Servername: strings.ToLower(servername),
		Type:       strings.ToLower(servertype),
		CPU:        icpu,
		RAM:        iram,
		Disk:       idisk,
	}

	p.Create()
}
