package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/High-la/snippetbox/internal/models"
)

// Change the signature of the home hander so it is defined as amethod against
// * application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is jut the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the new render helper.
	app.render(w, r, http.StatusOK, "home.tmpl.html", data)

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

	// And do the same thing again here...
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Use the new render helper.
	app.render(w, r, http.StatusOK, "view.tmpl.html", data)

}

// Change the signature of the snippetCreate handler so it is defined as a method
// against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

// Change the signature of the snippetCreatePost handler so it is defined as a method
// agains * application.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// First we call r.ParseForm() which adds any data in post request bodies
	// to the r.PostForm map. This also workd in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClentError() helper to
	// send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// The r.PostForm.Get() method always returns the form data as a *string*
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an integer. So we need to manually convert
	// the form data to an integer using strconv.Atoi(), and we send a 400 Bad
	// Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Initialize a map to hold any validation errors for the fields.
	fieldsErrors := make(map[string]string)

	// Check that the title value is not blank and is not more than 100
	// chars long. If it fails either of those checks, add a message to the
	// errors map using the field nae as the key.
	if strings.TrimSpace(title) == "" {
		fieldsErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldsErrors["title"] = "This field cannot be more than 100 characters long "
	}

	// Check that the Content value isn't blank.
	if strings.TrimSpace(content) == "" {
		fieldsErrors["content"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(content) > 250 {
		fieldsErrors["title"] = "This field cannot be more than 250 characters long "
	}

	// Check the expires value matches one of the permitted
	// values (1, 7, 365)
	if expires != 1 && expires != 7 && expires != 365 {
		fieldsErrors["expires"] = "This field must be 1, 7 or 365"
	}

	// If there are any errors, dump them in a plain text HTTP response and
	// return from the handler.
	if len(fieldsErrors) > 0 {
		fmt.Fprint(w, fieldsErrors)
		return
	}

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID for the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
