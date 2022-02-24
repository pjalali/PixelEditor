package routes

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	"pjalali.github.io/pixeleditor/internal/pkg/imageUtils"
)

type EditorTemplateData struct {
	OriginalImage string
	ImgData       string
	Height        int
	Width         int
	Threads       int
	Red           int
	Blue          int
	Green         int
	Contrast      int
	Hue           int
	Saturation    int
	Light         int
	TimeRGB       string
	TimeHSL       string
	TimeTotal     string
}

func Capture(w http.ResponseWriter, r *http.Request) {
	log.Println("WITHIN CAPUTRE")
	editorTemplateFileLocation := filepath.Join("templates", "editor.html")
	log.Println(editorTemplateFileLocation)
	tmpl, _ := template.ParseFiles(editorTemplateFileLocation)
	r.ParseForm()

	edt := EditorTemplateData{
		OriginalImage: r.FormValue("imgData"),
		ImgData:       r.FormValue("imgData"),
		Height:        stringToInt(r.FormValue("hData")),
		Width:         stringToInt(r.FormValue("wData")),
		Threads:       1,
		Red:           0,
		Blue:          0,
		Green:         0,
		Contrast:      0,
		Hue:           0,
		Saturation:    0,
		Light:         0,
		TimeRGB:       "",
		TimeHSL:       "",
		TimeTotal:     "",
	}

	tmpl.Execute(w, edt)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	log.Println("Within the editor")
	editorTemplateFileLocation := filepath.Join("templates", "editor.html")
	tmpl, _ := template.ParseFiles(editorTemplateFileLocation)
	r.ParseForm()

	edt := httpRequestToEditorTemplateData(r)

	base64ByteArray, err := base64.StdEncoding.DecodeString(edt.OriginalImage[22:])
	if err != nil {
		log.Println("Cannot decode base64 string")
		panic(err)
	}

	reader := bytes.NewReader(base64ByteArray)
	base64Image, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	base64ImageBounds := base64Image.Bounds()
	base64ImageRGBA := image.NewRGBA(image.Rect(0, 0, base64ImageBounds.Dx(), base64ImageBounds.Dy()))
	draw.Draw(base64ImageRGBA, base64ImageRGBA.Bounds(), base64Image, base64ImageBounds.Min, draw.Src)

	if edt.Threads < 1 || edt.Threads > edt.Height {
		panic("Invalid number of threads.")
	}

	if edt.Red != 0 || edt.Green != 0 || edt.Blue != 0 {
		imageUtils.ModifyRGBParallel(base64ImageRGBA, edt.Red, edt.Green, edt.Blue, edt.Threads)
	}

	if edt.Contrast != 0 {
		imageUtils.ModifyContrastParallel(base64ImageRGBA, edt.Contrast, edt.Threads)
	}

	if edt.Hue != 0 || edt.Saturation != 0 || edt.Light != 0 {
		imageUtils.ImageHSLModifications(base64ImageRGBA, edt.Hue, edt.Saturation, edt.Light)
	}

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, base64ImageRGBA)
	if err != nil {
		log.Println("Unable to encode image to png base64.")
		panic(err)
	}

	modifiedImgBase64Data := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buffer.Bytes())

	edt.ImgData = modifiedImgBase64Data

	tmpl.Execute(w, edt)
}

func stringToInt(value string) int {
	intVal, err := strconv.Atoi(value)
	if err != nil {
		log.Println("Unable to convert string to int.")
		panic(err)
	}
	return intVal
}

func httpRequestToEditorTemplateData(r *http.Request) EditorTemplateData {

	edt := EditorTemplateData{
		OriginalImage: r.FormValue("originalImage"),
		Height:        stringToInt(r.FormValue("hData")),
		Width:         stringToInt(r.FormValue("wData")),
		Threads:       stringToInt(r.FormValue("threads")),
		Red:           stringToInt(r.FormValue("rOffset")),
		Green:         stringToInt(r.FormValue("gOffset")),
		Blue:          stringToInt(r.FormValue("bOffset")),
		Contrast:      stringToInt(r.FormValue("contrast")),
		Hue:           stringToInt(r.FormValue("hue")),
		Saturation:    stringToInt(r.FormValue("sat")),
		Light:         stringToInt(r.FormValue("light")),
		TimeRGB:       "",
		TimeHSL:       "",
		TimeTotal:     "",
	}
	return edt
}
