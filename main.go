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

	"github.com/gin-gonic/gin"
)

//go:embed static/* templates/*
var f embed.FS

func generateEmailPrefix() string {
	hasher := md5.New()
	fmt.Fprint(hasher, time.Now().UTC().Unix())
	fmt.Fprint(hasher, rand.Intn(1000))
	return hex.EncodeToString(hasher.Sum(nil))[:16]
}

func main() {
	gin.DisableConsoleColor()
	router := gin.Default()

	router.SetHTMLTemplate(template.Must(template.New("").ParseFS(f, "templates/*.gohtml")))

	static, _ := fs.Sub(f, "static")
	router.StaticFS("/static", http.FS(static))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gohtml", gin.H{
			"emailPrefix": generateEmailPrefix(),
		})
	})

	err := router.Run()
	if err != nil {
		panic(err)
	}
}
