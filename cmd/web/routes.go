package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.

// Update the signature for the routes() method so that it returns a
// http.handler instead of *http.ServeMux
func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" dir.
	// Note that the path given to the http.Dir function is relative to the project
	// dir root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Swap the route declarations to use the application struct's methods as the
	// handler function.
	// Register other application routes as normal.

	// Update these routes to use the new dynamic middleware chain followed by
	// the appropriate handler function. Note that because the alice ThenFunc()
	// method returns a http.Handler (rather than a http.HandlerFunc) we also
	// need to switch to registering the route using the mux.Handle() method.
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// Pass the servemux as the 'next' parameter to the commonHeaders middleware
	// Because commonHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.

	// Wrap the existing chain with the logRequest middleware.
	// Wrap the existing chain with the recoverPanic middleware.

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request of our app receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Retrun the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
}
