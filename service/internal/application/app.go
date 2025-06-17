package application

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//go:embed init.sql
var initSQL string

type App struct {
	Config   Config
	LogFile  *os.File
	Logger   *slog.Logger
	Router   *mux.Router
	Database *sql.DB
}

func NewApp(config Config) (*App, error) {

	// Create data directory if it doesn't exist
	err := os.MkdirAll("data", 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// TODO: Maybe replace with postgres in production
	db, err := sql.Open("sqlite3", "data/app.db")
	if err != nil {
		panic(err)
	}

	// Check if database is initialized by looking for a known table
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' LIMIT 1").Scan(&tableName)
	if err == sql.ErrNoRows {
		// Database is not initialized, run init SQL
		_, err = db.Exec(initSQL)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}
	} else if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to check database initialization: %w", err)
	}

	// Open the log file in append mode, create if it doesn't exist
	logFile, err := os.OpenFile(config.Logging.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Determine the log level from the config
	var logLevel slog.Level
	switch config.Logging.LogLevel {
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
	switch config.Logging.LogFormat {
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
	if config.Server.Host != "" {
		router.Host(config.Server.Host)
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
		AllowedOrigins: []string{
			"http://localhost:3000", // React dev server
			"http://localhost:5173", // Vite dev server
			"http://127.0.0.1:3000", // Alternative localhost
			"http://127.0.0.1:5173", // Alternative localhost
			"t3.hinkolas.de",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(app.Router)

	server := &http.Server{
		Handler: handler,
		Addr:    ":3141",
	}

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
