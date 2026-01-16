package models

import (
	"database/sql"
	"time"
)

// Remember: The internal directory is being used to hold ancillary non-application-
// specific code, which could potentially be reused. A database model which could be
// used by other applications in the future (like a command line interface application) fits
// the bill here.

// Define a Snippet type to hold the data for an individual snippet. Notice
// how fields of the struct correspond to the fields in mysql snippets
// table ?
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into database.
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {

	return 0, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {

	return Snippet{}, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {

	return nil, nil
}
