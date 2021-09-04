package main

import (
	"encoding/base64"
	"image/png"
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

type EditorTemplateData struct {
	OriginalImage    string
	ImgData  		 string
	Height   		 string
	Width   		 string
	Threads    		 string
	Red   		     string
	Blue       		 string
	Green      		 string
	Contrast  		 string
	Hue       		 string
	Saturation 		 string
	Light     	 	 string
	TimeRGB    		 string
	TimeHSL    		 string
	TimeTotal  		 string
}

const PORT = ":5000"

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/capture", capture)
	http.HandleFunc("/edit", edit)

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

func edit(w http.ResponseWriter, r *http.Request) {
	log.Println("Within the editor")
	editorTemplateFileLocation := filepath.Join("templates", "editor.html")
	tmpl, _ := template.ParseFiles(editorTemplateFileLocation)
	r.ParseForm()
	originalImgDataB64 := r.FormValue("originalImage")
	width := r.FormValue("wData")
	height := r.FormValue("hData")

	unbased, err := base64.StdEncoding.DecodeString(originalImgDataB64[22:])
	if err != nil {
		log.Println("Cannot decode base64 string")
	}

	log.Println(unbased)

	reader := bytes.NewReader(unbased)
	image, err := png.Decode(reader)
	if err != nil {
		log.Println(err)
		log.Println("Bad image")
	}

	// perform calculations here

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, image)
	
	modifiedImgBase64Data := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buffer.Bytes())

	edt := EditorTemplateData{
		OriginalImage:	originalImgDataB64,
		ImgData:	    modifiedImgBase64Data,
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
