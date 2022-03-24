# Pixel Editor

## About

Originally done as the final project for Greg Baker's CMPT 383 class at Simon Fraser University, this program provides the means of editing photos that's both fun, efficient, and convenient. Users can take a photo of themselves using their webcam (permissions must be granted), or upload their own photos. **Note: to use the webcam features, one must have a working webcam.** The project is somewhat reminiscent of an application like Photo Booth that is provided on Apple computers, but the level of execution is much more on par with an undergraduate project rather than an application developed by a company worth trillions of dollars.

This program has been modified since its submission - mainly to get rid of the Python web server in favor of Go. To view the program in a closer state to what was handed in, take a look at the legacy branch.

## Languages Used

This project uses JavaScript and Go. JavaScript is used for the client-side code that manipulates the web pages, manages the webcam, validates user inputs, and sends POST/GET requests to the server as necessary. Go is used as the web server,and for editing the images. Initially, the Python using the Flask library was used for the web server.

## How to run

1. Run `go run main.go`
