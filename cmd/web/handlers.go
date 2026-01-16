package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/High-la/snippetbox/internal/models"
)

// Change the signature of the home hander so it is defined as amethod against
// * application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.

	// Include the navigation partial in the template files.
	// files := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.htmrowsl",
	// 	"./ui/html/pages/home.tmpl.html",
	// }

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message, use
	// the http.Error() function to send an Internal Server Error response to the
	// user, and then return from the  handler so no subsequent code is executed.

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we use ... to pass the contents
	// of the files slice as variadic arguments.
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// Because the home handler is now a method agains the application
	// struct it can access its fields, including the structured logger. we
	// use this to create a log entry at Error level containing the error
	// message, also including the request method and URI as attributes to
	// assist with debugging.
	// app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	// http.Error(w, "Internal Server Error", http.StatusInternalServerError)

	// use the serverErrorError() helper
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// Then we use the Execute() menthod on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil.

	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// And we also need to update the code here to use the structured logger too.
	// app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	// http.Error(w, "Internal Server Error", http.StatusInternalServerError)

	//  Use the serverError() helper
	// app.serverError(w, r, err)
	// }

	// w.Write([]byte("Hello from Snippetbox"))
}

// change the signature of the snippetView handler so it is defined as a method
// against * application
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the SnippetModel's Get() method to retrieve the data for a
	// specific record based on its ID. if no matching record is found,
	// return a 404 not found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Write the snippet data as a plain-text to HTTP response body
	fmt.Fprintf(w, "%+v", snippet)
}

// Change the signature of the snippetCreate handler so it is defined as a method
// against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// Change the signature of the snippetCreatePost handler so it is defined as a method
// agains * application.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// Create some variables holding dummy data. will be deleted
	// later during build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID for the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("Save a new snippet..."))
}
