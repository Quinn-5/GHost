package main

import (
	"html/template"
	"net/http"

	"github.com/Quinn-5/GHost/ghost"
	"github.com/Quinn-5/GHost/ghost/servconf"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/resource"
)

func createHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		servername := ctx.PostForm("servername")
		servertype := ctx.PostForm("servertype")
		var cpu resource.Quantity
		if n, err := resource.ParseQuantity(ctx.PostForm("cpu")); err == nil {
			cpu = n
		}
		var ram resource.Quantity
		if n, err := resource.ParseQuantity(ctx.PostForm("ram") + "Gi"); err == nil {
			ram = n
		}
		var disk resource.Quantity
		if n, err := resource.ParseQuantity(ctx.PostForm("disk") + "Gi"); err == nil {
			disk = n
		}

		conf := &servconf.ServerConfig{
			Username:   username,
			Servername: servername,
			Type:       servertype,
			CPU:        cpu,
			RAM:        ram,
			Disk:       disk,
		}

		err := ghost.Create(conf)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.SetCookie("username", conf.Username, 30, "/", "localhost", false, true)
		ctx.SetCookie("servername", conf.Servername, 30, "/", "localhost", false, true)

		ctx.Redirect(http.StatusFound, "/success/")
	}
}

func resultHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var username string
		if cookie, err := ctx.Cookie("username"); err == nil {
			println(cookie)
			username = cookie
		}
		var servername string
		if cookie, err := ctx.Cookie("servername"); err == nil {
			println(cookie)
			servername = cookie
		}

		conf := &servconf.ServerConfig{
			Username:   username,
			Servername: servername,
		}

		ghost.GetAddress(conf)

		ctx.HTML(http.StatusOK, "result.html", conf)
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

	route := gin.Default()
	route.Static("/static", "./static")
	route.StaticFile("/favicon.ico", "./resources/favicon.ico")
	route.StaticFile("/navbar.html", "./tmpl/navbar.html")
	route.LoadHTMLGlob("tmpl/*")

	route.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	route.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.html", gin.H{})
	})

	route.POST("/create", createHandler())

	route.GET("/success", resultHandler())

	route.Run(":8000")
}
