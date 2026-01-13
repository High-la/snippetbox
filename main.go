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

func main() {

	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL patter.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

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
