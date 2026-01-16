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

	// Write the SQL stmt we want to execute. it's splitted to two lines
	// for readability.
	stmt := `INSERT INTO snippets (title, content, created, expires)
			VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL stmt, followed by
	// values for the placeholder params. This method returns a sql.Result type which contains some
	// basic information bout what happened when the was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result to get the ID of newly
	// inserted record in the snippets table.
	// 	Important: Not all drivers and databases support the LastInsertId() and
	// RowsAffected() methods. For example, LastInsertId() is not supported by
	// PostgreSQL. So if you’re planning on using these methods it’s important to check the
	// documentation for your particular driver first.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so convert it to int type
	// before returning.
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {

	return Snippet{}, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {

	return nil, nil
}
