package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", homepageHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("unable to listen and serve:", err)
	}
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("index.html")
	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}