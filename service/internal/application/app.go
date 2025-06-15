package application

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Config struct {
	// Server
	Host         string // Hostname of the application
	ReadTimeout  int    // Time a request must take at most in seconds
	WriteTimeout int    // Time a response must take at most in seconds
	// Logging
	LogFilePath string // Path of the log file used for structured logging
	LogLevel    string // Set the logging level ("debug", "info", "warn", "error") (default "info")
	LogFormat   string // Format of the structured logger (text, json) (default: json)
	LogVerbose  bool   // Output log messages to stdout in addition to the log file
}

type App struct {
	Config   Config
	LogFile  *os.File
	Logger   *slog.Logger
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		app.Logger.Debug("Request received",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"query", r.URL.RawQuery,
		)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		app.Logger.Debug("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", duration.String(),
		)
	})
}

func NewApp(config Config) (*App, error) {

	// TODO: Maybe replace with postgres in production
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}

	// Open the log file in append mode, create if it doesn't exist
	logFile, err := os.OpenFile(config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Determine the log level from the config
	var logLevel slog.Level
	switch config.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Create a log handler with the provided format, output and level
	var logHandler slog.Handler
	switch config.LogFormat {
	case "json":
		logHandler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: logLevel, // Set the minimum log level
		})
	case "text":
		logHandler = slog.NewTextHandler(logFile, &slog.HandlerOptions{
			Level: logLevel, // Set the minimum log level
		})
	default:
		logHandler = slog.NewTextHandler(logFile, &slog.HandlerOptions{
			Level: logLevel, // Set the minimum log level
		})
	}

	logger := slog.New(logHandler)

	// Create new Router
	router := mux.NewRouter()

	// Set the valid Host if provided in the config
	if config.Host != "" {
		router.Host(config.Host)
	}

	return &App{
		Config:   config,
		LogFile:  logFile,
		Logger:   logger,
		Router:   router,
		Database: db,
	}, nil

}

func (app *App) Start() error {

	fmt.Println("Starting app...")

	// Add logging middleware to router
	app.Router.Use(app.loggingMiddleware)

	// Setup CORS middleware
	// TODO: use proper cors settings in production
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(app.Router)

	server := &http.Server{
		Handler: handler,
		Addr:    ":3141",
	}

	// Set ReadTimeout for server
	if app.Config.WriteTimeout > 0 {
		server.WriteTimeout = time.Duration(app.Config.WriteTimeout) * time.Second
	}

	// Set WriteTimeout for server
	if app.Config.ReadTimeout > 0 {
		server.WriteTimeout = time.Duration(app.Config.ReadTimeout) * time.Second
	}

	// TODO: connect to database

	return server.ListenAndServe()

}

func (app *App) Close() error {

	var err error

	err = app.LogFile.Close()
	if err != nil {
		return fmt.Errorf("error closing logfile: %w", err)
	}

	err = app.Database.Close()
	if err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}

	return err

}
