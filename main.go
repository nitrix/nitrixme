package main

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//go:embed static/* templates/* .well-known/*
var f embed.FS

func main() {
	gin.DisableConsoleColor()

	router := gin.New()
	router.Use(gin.Recovery())

	router.SetHTMLTemplate(template.Must(template.New("").ParseFS(f, "templates/*.gohtml")))

	static, _ := fs.Sub(f, "static")
	router.StaticFS("/static", http.FS(static))

	openpgpkey, _ := fs.Sub(f, ".well-known/openpgpkey")
	router.Group("/.well-known/openpgpkey", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Access-Control-Allow-Origin", "*")
	}).StaticFS("/", http.FS(openpgpkey))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gohtml", gin.H{})
	})

	// Support HTTP/2 over clear-text (h2c) for Cloud Run.
	if useh2c := os.Getenv("USE_H2C"); useh2c != "" {
		router.UseH2C = true
	}

	var err error
	if gin.Mode() == gin.DebugMode {
		err = router.Run("localhost:8080")
	} else {
		err = router.Run()
	}
	if err != nil {
		panic(err)
	}
}
