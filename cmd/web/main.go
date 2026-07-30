package web

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/slashtechno/generate-ddg/assets"
	"github.com/slashtechno/generate-ddg/internal"
)

// The application struct holds the dependencies needed for our handlers,
// including a htmlRenderer type.
// auth stuff is defined in auth.go
type application struct {
	logger    *slog.Logger
	html      *htmlRenderer
	jwtSecret string
	password  string
}

func Main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize a new htmlRenderer, parsing the base template and all partial
	// templates from assets/html into the shared template set.
	htmlRenderer, err := newHTMLRenderer(assets.HTMLFiles, "base.tmpl", "partials/*.tmpl")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	jwtSecret := internal.SecretViper.GetString("jwt-secret")
	if jwtSecret == "" {
		logger.Error("jwt-secret is not set; run \"generate-ddg config\" or set it in the secrets file/environment")
		os.Exit(1)
	}

	password := internal.SecretViper.GetString("web-password")
	if password == "" {
		logger.Error("web-password is not set; run \"generate-ddg config\" or set it in the secrets file/environment")
		os.Exit(1)
	}

	// Include the htmlRenderer in the application struct.
	app := &application{
		logger:    logger,
		html:      htmlRenderer,
		jwtSecret: jwtSecret,
		password:  password,
	}

	// Create a file server that serves the files from assets/static.
	fileserver := http.FileServerFS(assets.StaticFiles)

	// Protected routes: everything registered on protectedMux requires auth.
	// Paths here are relative to "/app" since http.StripPrefix removes it
	// before protectedMux ever sees the request.
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /gopher", app.gopher)
	protectedMux.HandleFunc("GET /logout", app.logout)

	// Register the application routes.
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("POST /login", app.login)
	mux.Handle("/app/", app.requireAuth(http.StripPrefix("/app", protectedMux)))

	// Start the HTTP server.
	logger.Info("starting server", "port", 5051)
	err = http.ListenAndServe(":5051", app.attachAuthStatus((mux)))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
