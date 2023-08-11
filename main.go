package main

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//go:embed static
var staticFiles embed.FS

var templates *template.Template

func main() {
	staticFS := http.FS(staticFiles)

	http.Handle("/static/", http.FileServer(staticFS))
	http.HandleFunc("/", homepageHandler)

	var err error

	templates, err = template.ParseFS(staticFiles, "static/*.gohtml")
	if err != nil {
		log.Fatalln(err)
	}

	port := 8080

	if envPort := os.Getenv("PORT"); envPort != "" {
		if n, err := strconv.Atoi(envPort); err == nil {
			port = n
		}
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalln("unable to listen and serve:", err)
	}
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	hasher := md5.New()
	fmt.Fprint(hasher, time.Now().UTC().Unix())
	fmt.Fprint(hasher, rand.Intn(1000))
	emailPrefix := hex.EncodeToString(hasher.Sum(nil))[:16]

	err := templates.ExecuteTemplate(w, "index.gohtml", struct {
		EmailPrefix string
	}{
		EmailPrefix: emailPrefix,
	})

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}
