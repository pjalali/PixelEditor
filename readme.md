# Pixel Editor

## About

Originally done as the final project for Greg Baker's CMPT 383 class at Simon Fraser University, this polygot program provides the means of editing photos that's both fun, efficient, and convenient. Users can take a photo of themselves using their webcam (permissions must be granted), or upload their own photos. **Note: to use the webcam features, one must have a working webcam.** The project is somewhat reminiscent of an application like Photo Booth that is provided on Apple computers, but the level of execution is much more on par with an undergraduate project rather than an application developed by a company worth trillions of dollars.

## Languages Used

This project used JavaScript, Python, and Go. JavaScript was used for the client-side code that manipulated the web pages, handled the managing the webcam and file uploads, validating user inputs, and sent POST/GET requests to the server as necessary. Python was used as the web server that utilized the Flask library, and Go was used for actually editing the images.

## Communication Methods

The client-side JavaScript code along with the HTML would send requests to the Python REST web server which would in turn send a response, and the Python server would run the Go code as a compiled executable.

After navigating to the webpage, the user would take/upload their photo and the client-side JavaScript would then send the image data as a POST request to the Python web server. The Python web server would then save that image locally, and respond to the POST request with a new webpage containing the image, and slider inputs for modifying the image. Once the user decided on a configuration using the sliders, they would submit and another POST request would be sent to the the Python web server with the values to modify. The web server would then execute the compiled Go binary giving it the values the user entered, and the Go program would write the edited image to a file. Once the Go program finished, the Python web server would send the modified image (along with the input sliders for further editing) back to the user as a response to the POST request.

## How to run

1. Ensure you have `go`, `python`, `python-flask`, `python-pil`, and `make` installed
1. Clone repository and `cd` into it
1. Run `export FLASK_APP=app.py`
1. Run `make`
1. Navigate to `127.0.0.1:5000`
1. Have fun!
