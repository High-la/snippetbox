package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	// Use the Header().Add() method to add a 'Server: Go' header to the
	// response header map. The first parameter is the header name, and
	// the second parameter is the header value.
	w.Header().Add("Server", "Go")

	// Manipulating the header map
	// Set a new cache-control header. If an existing "Cache-Control" header exists
	// it will be overwritten.
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	// In contrast, the Add() method appends a new "Cache-Control" header and can
	// be called multiple times.
	w.Header().Add("Cache-Control", "public")
	w.Header().Add("Cache-Control", "max-age=31536000")
	// Delete all values for the "Cache-Control" header.
	w.Header().Del("Cache-Control")
	// Retrieve the first value for the "Cache-Control" header.
	w.Header().Get("Cache-Control")
	// Retrieve a slice of all values for the "Cache-Control" header.
	w.Header().Values("Cache-Control")

	// header, set content type
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"name":"Haile"}`))

	w.Write([]byte("Hello from Snippetbox"))
}

// Add a snippetView handler function.
func snippetView(w http.ResponseWriter, r *http.Request) {

	// Extract the value of the id wildcard from the request using r.PathValue()
	// and try to convert it to an integer using the strconv.Atoi() function. If
	// it can't be converted to an integer, or the value is less than 1, we
	// return a 404 page not found response.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the fmt.Sprintf() function to interpolate the id value with a
	// message, then write it as the HTTP response.
	// msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)
	// w.Write([]byte(msg))

	// This will allow us to interpolate the wildcard id value in our
	// response body message and write the response in a single line, like so:
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// Add a snippetCreatePost handler function.
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// use the w.WriteHeader() method to send a 201 status code.
	w.WriteHeader(http.StatusCreated)

	// Then w.Write() method to write the response body as normal.
	w.Write([]byte("Save a new snippet..."))
}

func main() {

	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL patter.
	mux := http.NewServeMux()

	// "/" matches everything starts with "/" ends with any characters

	// 	Note: It’s only permitted to use {$} at the end of subtree path patterns (i.e. patterns
	// that end with a trailing slash). Route patterns that don’t have a trailing slash require a
	// match on the entire request path anyway, so it doesn’t make sense to include {$} at
	// the end and trying to do so will result in a runtime panic.

	// Restrict this route to exact matches on / only.
	mux.HandleFunc("GET /{$}", home)

	// Add the {id} wildcard segment
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	// Create the new route, which is restricted to POST request only.
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// Print a log message to say that server is starting.
	log.Print("Starting server on :4000")

	// Use the http.ListenAndServe() function to start a new web server. we pass in
	// two parameters: the TCP network address on(in this case ":4000")
	// and the servemux we just created. if http.ListenAndServer() returns an error
	//  we use log.Fatal() function to log the error message and exit.
	// Note that any error returned by http.ListenAndServer() is always non-nil
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
