package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	mux := http.NewServeMux()

	// os.Getenv() only reads from already setted system environment variables.
	// so we use the godotenv package to read the .env file and set the
	// environment variables before we call os.Getenv().
	godotenv.Load() // Load .env file
	addr := os.Getenv("SNIPPETBOX_ADDR")

	// Create a file server which serves files out of the "./ui/static" dir.
	// Note that the path given to the http.Dir function is relative to the project
	// dir root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Register other application routes as normal.
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Printf("starting server on %s", addr)
	err := http.ListenAndServe(addr, mux)
	log.Fatal(err)
}
