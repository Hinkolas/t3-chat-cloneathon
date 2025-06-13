package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application/auth"
	"github.com/spf13/cobra"
)

var DefaultConfig application.Config

func init() {
	rootCmd.AddCommand(startCmd)

	// Search Command Flags
	startCmd.Flags().BoolP("verbose", "v", false, "Output log messages to stdout in addition to the log file")
	startCmd.Flags().StringP("log-level", "l", "info", "Set the logging level (\"debug\", \"info\", \"warn\", \"error\", \"fatal\") (default \"info\")")

	DefaultConfig = application.Config{
		Host: ":3141",

		LogFilePath: "./app.log",
		LogLevel:    "debug",
		LogFormat:   "text",
	}
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the app with the provided configuration",
	Run: func(cmd *cobra.Command, args []string) {

		app, err := application.NewApp(DefaultConfig) // TODO: Replace with real config
		if err != nil {
			fmt.Printf("Error creating the application: %v\n", err)
		}
		defer app.Close()

		app.Router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
			app.Logger.Info("Application is live and well!")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Why are you checking? Don't you trust me?"})
		})

		authService, err := auth.NewService(app)
		if err != nil {
			fmt.Printf("Error initializing auth service: %v\n", err)
		}

		authService.Handle(app.Router)

		if err = app.Start(); err != nil {
			fmt.Printf("Error starting the application: %v\n", err)
		}

	},
}
