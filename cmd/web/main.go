package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"

	"github.com/High-la/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"

	"github.com/go-playground/form/v4"
)

type application struct {
	logger         *slog.Logger
	snippets       models.SnippetModelInterface // Use our new interface type.
	users          models.UserModelInterface    // Use our new interface type.
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	// Load environment variables from .env
	godotenv.Load()

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// logs caller file and line number
		AddSource: true,
		// u can also add custom level names or disable level display
		Level: slog.LevelDebug,
	}))

	// --------------------
	// Database
	// --------------------
	dsn := os.Getenv("SNIPPETBOX_DB_DSN")
	db, err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// --------------------
	// Templates & Forms
	// --------------------
	// Intitialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()

	// --------------------
	// Sessions
	// --------------------
	// Use the scs.New() function to initialize a new session manager. then we
	// configure it to use our MySQL database as the session store, and set a
	// lifetime of 12 hours (so that sessions auto expires 12 hours)
	// after first being created).
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// --------------------
	// App
	// --------------------
	// Initialize a new instance of our application struct, containinig the
	// dependencies..
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// --------------------
	// TLS config
	// --------------------
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
	}

	// --------------------
	// HTTP Server
	// --------------------
	addr := fmt.Sprintf(
		"%s:%s",
		os.Getenv("SNIPPETBOX_BIND_ADDR"),
		os.Getenv("SNIPPETBOX_BIND_PORT"),
	)

	// 	Create an http.Server (important)

	// Don’t use http.ListenAndServe
	// Always use http.Server
	// Why?
	// http.Server exposes Shutdown() → this is the magic

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before.

	// Set the server's TLSConfig field to use the tlsConfig we just
	// created
	srv := &http.Server{
		Addr:    addr,
		Handler: app.routes(),

		// Create a *log.Logger from our structured logger handler, which writes
		// log entries at Error level, and assign it to the Error log field. IF
		// u would prefer to log the server errors at Warn level instead, u
		// could pass slog.LevelWarn aas the final parameter.
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// --------------------
	// Graceful shutdown
	// --------------------
	// App exits cleanly when Docker / VPS stops it.
	// clean shutdown logs → no panic
	// 	On shutdown (SIGINT, SIGTERM):

	// 1. Listen for OS signals
	// 2. Stop accepting new HTTP requests
	// 3. Let in-flight requests finish
	// 4. Close DB connections cleanly

	// Listen for SIGINT / SIGTERM
	// Create a channel to receive OS signals.
	shutdownCh := make(chan os.Signal, 1)

	// Notifyon SIGINT (Ctrl+c) AND SIGTERM(Docker, systemd).
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	// Start Server in a Goroutine
	// The main goroutine must stay free to listen for shutdown signals.
	go func() {
		logger.Info("starting server", "addr", addr, "env=", os.Getenv("SNIPPETBOX_ENV"))

		err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
		if err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
		}
	}()

	// Wait for shutdown signal
	<-shutdownCh
	app.logger.Info("shutdown signal received")

	// Graceful shutdown timeout (env-based)
	shutdownTimeout := 5 * time.Second
	if v := os.Getenv("SNIPPETBOX_SHUTDOWN_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			shutdownTimeout = d
		}
	}

	// Stop Accepting New Requests (Gracefully)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
	}

	logger.Info("server stopped cleanly")

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
