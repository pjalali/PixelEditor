package main

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

type EditorTemplateData struct {
	ImgData    string
	Height     string
	Width      string
	Threads    string
	Red        string
	Blue       string
	Green      string
	Contrast   string
	Hue        string
	Saturation string
	Light      string
	TimeRGB    string
	TimeHSL    string
	TimeTotal  string
}

const PORT = ":5000"

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/capture", capture)

	log.Println("Listening on " + PORT)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Println(err)
	}
}

func capture(w http.ResponseWriter, r *http.Request) {
	editorTemplateFileLocation := filepath.Join("templates", "editor.html")
	tmpl, _ := template.ParseFiles(editorTemplateFileLocation)
	r.ParseForm()
	imgDataB64 := r.FormValue("imgData")
	width := r.FormValue("wData")
	height := r.FormValue("hData")

	edt := EditorTemplateData{
		OriginalImage:	imgDataB64,
		ImgData:	    imgDataB64,
		Height:     	height,
		Width:	 	    width,
		Threads:	    "1",
		Red: 	        "0",
		Blue: 	        "0",
		Green:  	    "0",
		Contrast:	    "0",
		Hue:   		    "0",
		Saturation: 	"0",
		Light:      	"0",
		TimeRGB:    	"",
		TimeHSL:    	"",
		TimeTotal:  	"",
	}

	tmpl.Execute(w, edt)

}
