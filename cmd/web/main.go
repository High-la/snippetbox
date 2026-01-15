package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// os.Getenv() only reads from already setted system environment variables.
	// so we use the godotenv package to read the .env file and set the
	// environment variables before we call os.Getenv().
	godotenv.Load() // Load .env file
	addr := os.Getenv("SNIPPETBOX_ADDR")

	mux := http.NewServeMux()

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

	// Use the Info() method to log the starting server message at Info severity
	// (along with the listen address as an attribute).
	logger.Info("starting server", "addr", addr)

	err := http.ListenAndServe(addr, mux)

	// And we also use the Error() method to log any error message returned by
	// http.ListenAndServer() at Error severity (with no additional attributes),
	// and then call os.Exit(1) to terminate the application with exit code 1.
	logger.Error(err.Error())
	os.Exit(1)
}
