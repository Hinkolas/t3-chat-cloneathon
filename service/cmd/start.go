package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application/auth"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application/chat"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)

	// Start Command Flags
	startCmd.Flags().BoolP("verbose", "v", false, "Output log messages to stdout in addition to the log file")
	startCmd.Flags().StringP("log-level", "l", "info", "Set the logging level (\"debug\", \"info\", \"warn\", \"error\", \"fatal\") (default \"info\")")
	startCmd.Flags().StringP("config", "c", "config.yaml", "Path to config file")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the app with the provided configuration",
	Run: func(cmd *cobra.Command, args []string) {

		cfgFile := cmd.Flag("config").Value.String()

		cfg, err := application.LoadConfig(cfgFile) // Load config from file
		if err != nil {
			fmt.Printf("Error loading config from file: %v\n", err)
			return
		}

		// Create new application with loaded config
		app, err := application.NewApp(*cfg)
		if err != nil {
			fmt.Printf("Error creating the application: %v\n", err)
			return
		}
		defer app.Close()

		// Register basic endpoint to get application status
		app.Router.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
			app.Logger.Info("Application is live and well!")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Why are you checking? Don't you trust me?"})
		})

		authService, err := auth.NewService(app)
		if err != nil {
			fmt.Printf("Error initializing auth service: %v\n", err)
		}

		authService.Handle(app.Router)

		// T3 Chat
		chatService, err := chat.NewService(app)
		if err != nil {
			fmt.Printf("Error initializing auth service: %v\n", err)
		}

		// Add all models from the config file
		for key, model := range cfg.Models {
			// key is e.g. "claude-4-sonnet"
			chatService.AddModel(key, model)
		}

		chatService.Handle(app.Router)

		if err = app.Start(); err != nil {
			fmt.Printf("Error starting the application: %v\n", err)
		}

	},
}
