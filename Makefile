all: pixeleditorengine flask

pixeleditorengine:
	cd gorender; go build -o ../PixelEditorEngine main/main.go

flask:
	FLASK_APP=app.py; flask run