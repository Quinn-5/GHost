package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Quinn-5/GHost/ghost"
	"github.com/Quinn-5/GHost/ghost/configs/configstore"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func createRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	r.AddFromFiles("console", "templates/base.html", "templates/console.html")
	r.AddFromFiles("create", "templates/base.html", "templates/create.html")
	r.AddFromFiles("login", "templates/base.html", "templates/login.html")
	r.AddFromFiles("result", "templates/base.html", "templates/result.html")
	r.AddFromFiles("info", "templates/base.html", "templates/serverconsole/edit.html", "templates/serverconsole/info.html")
	r.AddFromFiles("settings", "templates/base.html", "templates/serverconsole/edit.html", "templates/serverconsole/settings.html")
	r.AddFromFiles("terminal", "templates/base.html", "templates/serverconsole/edit.html", "templates/serverconsole/terminal.html")
	return r
}

func cookieCheck(ctx *gin.Context) string {
	var username string
	if cookie, err := ctx.Cookie("username"); err != nil {
		ctx.Redirect(http.StatusFound, "/login")
	} else {
		username = cookie
	}
	return username
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

		ctx.HTML(http.StatusOK, "result", conf.Get())
	}
}

func shellHandler() gin.HandlerFunc {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return func(c *gin.Context) {
		username := cookieCheck(c)

		servername := c.Param("server")
		conf := configstore.New(username, servername)
		dataOut, dataIn := ghost.NewTerminal(conf.Get())

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.Writer.Write([]byte(fmt.Sprint("", err)))
			return
		}
		defer conn.Close()

		go func() {
			for {
				buf := make([]byte, 1024)
				dataOut.Read(buf)
				err = conn.WriteMessage(1, buf)
				if err != nil {
					log.Println("Error writing message: ", err)
					return
				}
			}
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message: ", err)
				return
			}
			dataIn.Write(append(message, byte('\n')))
		}
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
		username := cookieCheck(c)
		deployments := ghost.GetAllDeploymentsForUser(configstore.New(username, "").Get())
		c.HTML(http.StatusOK, "console", gin.H{
			"Servers": deployments,
		})
	})

	router.POST("/console", func(c *gin.Context) {
		username := cookieCheck(c)

		servername := c.PostForm("servername")
		conf := configstore.New(username, servername)

		action := strings.ToLower(c.PostForm("action"))
		switch action {
		case "edit":
			c.Redirect(http.StatusFound, "/console/"+servername)
			return
		case "delete":
			ghost.Delete(conf.Get())
		}

		c.Redirect(http.StatusFound, "/console")
	})

	router.GET("/console/:server", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/console/"+c.Param("server")+"/info")
	})

	router.GET("/console/:server/info", func(c *gin.Context) {
		username := cookieCheck(c)

		servername := c.Param("server")
		conf := configstore.New(username, servername)

		c.HTML(http.StatusOK, "info", conf.Get())
	})

	router.GET("/console/:server/settings", func(c *gin.Context) {
		username := cookieCheck(c)

		servername := c.Param("server")
		conf := configstore.New(username, servername)

		c.HTML(http.StatusOK, "settings", conf.Get())
	})

	router.GET("/console/:server/terminal", func(c *gin.Context) {
		username := cookieCheck(c)

		servername := c.Param("server")
		conf := configstore.New(username, servername)

		c.HTML(http.StatusOK, "terminal", conf.Get())
	})

	router.GET("/console/:server/terminal/shell", shellHandler())

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login", gin.H{})
	})

	router.POST("/login", func(c *gin.Context) {
		c.SetCookie("username", c.PostForm("username"), 3600, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, "/")
	})

	router.Run(":8000")
}
