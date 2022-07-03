package main

import (
	"net/http"

	"github.com/Quinn-5/GHost/ghost"
	"github.com/Quinn-5/GHost/ghost/servconf"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/resource"
)

func createRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	r.AddFromFiles("console", "templates/base.html", "templates/console.html")
	r.AddFromFiles("create", "templates/base.html", "templates/create.html")
	r.AddFromFiles("login", "templates/base.html", "templates/login.html")
	r.AddFromFiles("result", "templates/base.html", "templates/result.html")
	return r
}

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

		ctx.Redirect(http.StatusFound, "/success")
	}
}

func resultHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var username string
		if cookie, err := ctx.Cookie("username"); err == nil {
			username = cookie
		}
		var servername string
		if cookie, err := ctx.Cookie("servername"); err == nil {
			servername = cookie
		}

		conf := &servconf.ServerConfig{
			Username:   username,
			Servername: servername,
		}

		ghost.GetAddress(conf)

		ctx.HTML(http.StatusOK, "result", conf)
	}
}

func main() {
	router := gin.Default()
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.HTMLRender = createRenderer()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{})
	})

	router.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create", gin.H{})
	})

	router.POST("/create", createHandler())

	router.GET("/success", resultHandler())

	router.Run(":8000")
}
