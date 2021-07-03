package main

import (
	"log"
	"net/http"
)

const PORT = ":5000"

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Listening on " + PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Println(err)
	}
}
