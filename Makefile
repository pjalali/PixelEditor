all: pixeleditorengine flask

pixeleditorengine:
	cd gorender; GOPATH=`pwd` go build -o ../PixelEditorEngine src/main/main.go

flask:
	FLASK_APP=app.py; flask run