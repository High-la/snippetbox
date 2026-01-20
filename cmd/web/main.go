package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql" // New import
	//
	//	To use this model in our handlers we need to establish a new SnippetModel struct in our
	//
	// main() function and then inject it as a dependency via the application struct — just like we
	// have with our other dependencies.
	// Here’s how:
	// Import models package
	"github.com/High-la/snippetbox/internal/models"

	"github.com/go-playground/form/v4"
)

// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our hadlers.

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.

// Add a formDecoder field to hold a pointer to a form.Decoder instance.

// Add a new sessionManager field to the application struct.
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	// Define a new command-line flag for mysql DSN string.
	dsn := flag.String("dsn", "web:1234@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// logs caller file and line number
		AddSource: true,
		// u can also add custom level names or disable level display
		Level: slog.LevelDebug,
	}))

	// To keep the main() function tidy put the code for creating a connection
	// pool into the separate openDB() function below. we pass and pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// also defer a call to db.Close(), so that the connection pool is closed
	// before the main function exits.
	defer db.Close()

	// Intitialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()

	// Use the scs.New() function to initialize a new session manager. then we
	// configure it to use our MySQL database as the session store, and set a
	// lifetime of 12 hours (so that sessions auto expires 12 hours)
	// after first being created).
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Initialize a new instance of our application struct, containinig the
	// dependencies..

	// Add the session manager to our application dependecies.
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// os.Getenv() only reads from already setted system environment variables.
	// so we use the godotenv package to read the .env file and set the
	// environment variables before we call os.Getenv().
	godotenv.Load() // Load .env file
	addr := os.Getenv("SNIPPETBOX_ADDR")

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before.
	srv := &http.Server{
		Addr:    addr,
		Handler: app.routes(),

		// Create a *log.Logger from our structured logger handler, which writes
		// log entries at Error level, and assign it to the Error log field. IF
		// u would prefer to log the server errors at Warn level instead, u
		// could pass slog.LevelWarn aas the final parameter.
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Use the Info() method to log the starting server message at Info severity
	// (along with the listen address as an attribute).
	logger.Info("starting server", "addr", addr)

	// Call the ListenAndServer() method on our new http.Server struct to start
	// the server.
	err = srv.ListenAndServe()

	// And we also use the Error() method to log any error message returned by
	// http.ListenAndServer() at Error severity (with no additional attributes),
	// and then call os.Exit(1) to terminate the application with exit code 1.
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
