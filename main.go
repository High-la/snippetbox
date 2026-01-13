package main

import (
	"log"
	"net/http"
)

// Define a home handler function which writes a byte slice containing
// "Hello from snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

// Add a snippetView handler function.
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
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
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

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
