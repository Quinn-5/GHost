package main

import (
	"net/http"

	"github.com/Quinn-5/GHost/ghost"
	"github.com/Quinn-5/GHost/ghost/configs/configstore"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
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
		var username string
		if cookie, err := ctx.Cookie("username"); err == nil {
			username = cookie
		} else {
			ctx.Redirect(http.StatusFound, "/login")
		}
		servername := ctx.PostForm("servername")
		servertype := ctx.PostForm("servertype")
		cpu := ctx.PostForm("cpu")
		ram := ctx.PostForm("ram")
		disk := ctx.PostForm("disk")

		conf := configstore.New(username, servername)
		conf.SetType(servertype)
		conf.SetCPU(cpu)
		conf.SetRAM(ram)
		conf.SetDisk(disk)

		err := ghost.Create(conf.Get())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.SetCookie("servername", conf.Get().ServerName, 30, "/", "localhost", false, true)

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

		conf := configstore.New(username, servername)

		ghost.GetAddress(conf)

		ctx.HTML(http.StatusOK, "result", conf.Get())
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
		if _, err := c.Cookie("username"); err != nil {
			c.Redirect(http.StatusFound, "/login")
		}
		c.HTML(http.StatusOK, "create", gin.H{})
	})

	router.POST("/create", createHandler())

	router.GET("/success", resultHandler())

	router.GET("/console", func(c *gin.Context) {
		var username string
		if cookie, err := c.Cookie("username"); err != nil {
			c.Redirect(http.StatusFound, "/login")
		} else {
			username = cookie
		}
		deployments := ghost.GetAllDeploymentsForUser(configstore.New(username, "").Get())
		c.HTML(http.StatusOK, "console", gin.H{
			"Servers": deployments,
		})
	})

	router.POST("/console", func(c *gin.Context) {
		var username string
		if cookie, err := c.Cookie("username"); err != nil {
			c.Redirect(http.StatusFound, "/login")
		} else {
			username = cookie
		}

		servername := c.PostForm("servername")
		conf := configstore.New(username, servername)

		ghost.Delete(conf.Get())

		deployments := ghost.GetAllDeploymentsForUser(configstore.New(username, "").Get())
		c.HTML(http.StatusOK, "console", gin.H{
			"Servers": deployments,
		})
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login", gin.H{})
	})

	router.POST("/login", func(c *gin.Context) {
		c.SetCookie("username", c.PostForm("username"), 3600, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, "/")
	})

	router.Run(":8000")
}
