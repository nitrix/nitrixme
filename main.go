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

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/gin-gonic/gin"
	"go.opencensus.io/trace"
)

//go:embed static/* templates/*
var f embed.FS

func generateEmailPrefix() string {
	hasher := md5.New()
	fmt.Fprint(hasher, time.Now().UTC().Unix())
	fmt.Fprint(hasher, rand.Intn(1000))
	return hex.EncodeToString(hasher.Sum(nil))[:16]
}

func cloudTracing() {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})

	// Only begin tracing if we successfully created the exporter.
	if err == nil {
		trace.RegisterExporter(exporter)
	}
}

func init() {
	cloudTracing()
	gin.DisableConsoleColor()
}

func main() {
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
