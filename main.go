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
	servername := "sugma"
	servertype := "Minecraft"
	var cpu resource.Quantity
	if n, err := resource.ParseQuantity("1"); err == nil {
		cpu = n
	}
	var ram resource.Quantity
	if n, err := resource.ParseQuantity("2Gi"); err == nil {
		ram = n
	}
	var disk resource.Quantity
	if n, err := resource.ParseQuantity("1Gi"); err == nil {
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

	ghost.Create(p)

}
