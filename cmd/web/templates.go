package main

import (
	"html/template"
	"path/filepath"

	"github.com/High-la/snippetbox/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.

// Include a snippets field in the templateData struct.
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {

	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all file paths that
	// match the patter "./ui/html/pages/*.tmpl". This will essentially gives
	// us a slice of allthe filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop through the pages filepaths one by one.
	for _, page := range pages {

		// Extract the file name (like 'home.tmpl') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the files into a template set.
		ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Create a slice containing the filepaths for our base template, any
		// partials and the page.
		// files := []string{
		// 	"./ui/html/base.tmpl.html",
		// 	"./ui/html/partials/nav.tmpl.html",
		// 	page,
		// }

		// Add the template set to the map, using name of the page
		// (like 'home.tmpl') as the key

		// Add the template set to the map as normal...
		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}
