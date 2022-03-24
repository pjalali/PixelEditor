package main

import (
	"log"
	"net/http"
	"os"

	routes "pjalali.github.io/pixeleditor/internal/server"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		panic("$PORT must be set")
	}

	fs := http.FileServer(http.Dir("./internal/server/web/static"))
	http.Handle("/", fs)
	http.HandleFunc("/capture", routes.Capture)
	http.HandleFunc("/edit", routes.Edit)

	log.Println("Listening on " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println(err)
	}
}
