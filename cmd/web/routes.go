package main

import (
	"net/http"

	"github.com/High-la/snippetbox/ui"
	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.

// Update the signature for the routes() method so that it returns a
// http.handler instead of *http.ServeMux
func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	// Use the http.FileServerFS() function to create a HTTP handler which
	// serves the embedded files in ui.Files. It's important to note that our
	// static files are contained in the "static" folder of the ui.files
	// embedded filesystem. So, for eg, our CSS is located at
	// "static/css/main.css". This means that we no longer need to strip the
	// prefix from the request URL - any requests that start with /static/ can
	// can just be passed directly to the file server and the corresponding static
	// file will be served (so long as it exists).
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later

	// Unprotected application routes using the "dynamic" middleware chain.
	// Use the nosurf middleware on all our 'dynamic' routes.

	// Add the authenticate() middleware to the chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Swap the route declarations to use the application struct's methods as the
	// handler function.
	// Register other application routes as normal.

	// Update these routes to use the new dynamic middleware chain followed by
	// the appropriate handler function. Note that because the alice ThenFunc()
	// method returns a http.Handler (rather than a http.HandlerFunc) we also
	// need to switch to registering the route using the mux.Handle() method.
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected (authenticated-only) application routes, using a new 'protected'
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

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
