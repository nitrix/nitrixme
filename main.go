package main

import (
	"crypto/md5"
	"embed"
	_ "embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
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

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("unable to listen and serve:", err)
	}
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {

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
