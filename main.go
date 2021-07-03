package main

import (
	"fmt"
	"net/http"
)

const PORT = ":5000"

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("Listening on " + PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println("Error listening on " + PORT)
		fmt.Println(err)
	}
}
