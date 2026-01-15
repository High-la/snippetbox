package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.
type application struct {
	logger *slog.Logger
}

func main() {

	// Use the slog.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// logs caller file and line number
		AddSource: true,
		// u can also add custom level names or disable level display
		Level: slog.LevelDebug,
	}))

	// Initialize a new instance of our application struct, containinig the
	// dependencies (for now, the structured logger).
	app := &application{
		logger: logger,
	}

	// os.Getenv() only reads from already setted system environment variables.
	// so we use the godotenv package to read the .env file and set the
	// environment variables before we call os.Getenv().
	godotenv.Load() // Load .env file
	addr := os.Getenv("SNIPPETBOX_ADDR")

	// Use the Info() method to log the starting server message at Info severity
	// (along with the listen address as an attribute).
	logger.Info("starting server", "addr", addr)

	err := http.ListenAndServe(addr, app.routes())

	// And we also use the Error() method to log any error message returned by
	// http.ListenAndServer() at Error severity (with no additional attributes),
	// and then call os.Exit(1) to terminate the application with exit code 1.
	logger.Error(err.Error())
	os.Exit(1)
}
