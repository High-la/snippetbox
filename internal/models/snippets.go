package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (Snippet, error)
	Latest() ([]Snippet, error)
}

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

	// Write the SQL stmt we wanted to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > UTC_TIMESTAMP() AND id = ?`

	//  Use the QueryRow() method on the connection pool to execute the
	// SQL stmt, passing int the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a new zeroed Snippet struct.
	var s Snippet

	// Use row.Scan() to copy the values from each fields in sql.Row to the
	// corresponding field in the Snippet struct. Notice that arguments
	// to row.Scan are *pointers* to the place u want to copy the data into,
	// and the number of args must be exactly the same as the number of the
	// columns returned by ur stmmt.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function to check for that
		// error specifically, and return our own ErrNoRecord error
		// instead (we'll create this in a moment).

		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	// If everything went ok, then return the filled Snippet struct.
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {

	// Write the SQL stmt ...
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Use the Query() method on the conn pool to execute sql stmt
	// This returns a sql.Rows resultset containing the result of
	// our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// stmt should come *after* u check for an error from the Query()
	// method. Otherwise, if the Query() returns an error, u will get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs.
	var snippets []Snippet

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by
	// wor.Scan() method. If iteration over all the rows completes then the
	// resultset auto closes itself and fress-up the underlying
	// db conn.
	for rows.Next() {

		// Create a new zeroed Snippet struct.
		var s Snippet
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the args to row.Scan()
		// must be pointers to the place u want to copy the data into, and the
		// columns returned by your stmt.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)

	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice
	return snippets, nil
}
