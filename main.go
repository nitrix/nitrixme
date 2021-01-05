package main

import (
	_ "embed"
	"log"
	"net/http"
)

//go:embed index.html
var indexPage []byte

//go:embed nitrixme.jpg
var picture []byte

func main() {
	http.HandleFunc("/nitrixme.jpg", pictureHandler)
	http.HandleFunc("/", homepageHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("unable to listen and serve:", err)
	}
}

func pictureHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(picture)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(indexPage)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}