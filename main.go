package main

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/fs"
	"math/rand"
	"net/http"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/gin-gonic/gin"
)

//go:embed static/* templates/* .well-known/*
var f embed.FS

func generateEmailPrefix() string {
	hasher := md5.New()
	fmt.Fprint(hasher, time.Now().UTC().Unix())
	fmt.Fprint(hasher, rand.Intn(1000))
	return hex.EncodeToString(hasher.Sum(nil))[:16]
}

func main() {
	cfg := profiler.Config{
		Service:      "nitrixme",
		DebugLogging: true,
	}

	profiler.Start(cfg)

	gin.DisableConsoleColor()

	router := gin.Default()

	router.SetHTMLTemplate(template.Must(template.New("").ParseFS(f, "templates/*.gohtml")))

	static, _ := fs.Sub(f, "static")
	router.StaticFS("/static", http.FS(static))

	openpgpkey, _ := fs.Sub(f, ".well-known/openpgpkey")
	router.Group("/.well-known/openpgpkey", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Access-Control-Allow-Origin", "*")
	}).StaticFS("/", http.FS(openpgpkey))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gohtml", gin.H{
			"emailPrefix": generateEmailPrefix(),
		})
	})

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
