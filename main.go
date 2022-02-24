package main

import (
	"log"
	"net/http"

	routes "pjalali.github.io/pixeleditor/internal/server"
)

const PORT = ":5000"

func main() {
	fs := http.FileServer(http.Dir("./server/web/static"))
	http.Handle("/", fs)
	http.HandleFunc("/capture", routes.Capture)
	http.HandleFunc("/edit", routes.Edit)

	log.Println("Listening on " + PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Println(err)
	}
}
