package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application/auth"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application/chat"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm"
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

		// T3 Chat
		chatService, err := chat.NewService(app)
		if err != nil {
			fmt.Printf("Error initializing auth service: %v\n", err)
		}

		chatService.AddModel("claude-4-sonnet", llm.Model{
			Title:       "Claude 4 Sonnet",
			Description: "The latest model from Anthropic. Claude 4 Sonnet is a significant upgrade to Claude 3.7 Sonnet, delivering superior coding and reasoning while responding more precisely to your instructions.",
			Icon:        "anthropic",
			Name:        "claude-4-opus-20250514",
			Provider:    llm.Anthropic,
			Features: llm.ModelFeatures{
				HasVision: true,
				HasPDF:    true,
			},
			Flags: llm.ModelFlags{
				IsPremium:     true,
				IsRecommended: true,
			},
		})

		chatService.AddModel("qwen3", llm.Model{
			Title:       "Qwen3",
			Description: "An open source mixture-of-experts (MoE) language model developed by Alibaba Cloud, activating only 3 billion parameters out of a total of 30B. It comes in various sizes and is licenced under the Apache 2.0 license.",
			Icon:        "qwen",
			Name:        "qwen3:30b",
			Provider:    llm.Ollama,
			Features: llm.ModelFeatures{
				HasReasoning: true,
			},
			Flags: llm.ModelFlags{
				IsOpenSource: true,
				IsFree:       true,
			},
		})

		chatService.AddModel("gemini-2.0-flash", llm.Model{
			Title:       "Gemini 2.0 Flash",
			Description: "Google's flagship model, known for speed and accuracy (and also web search!). Not quite as smart as Claude 3.5 Sonnet, but WAY faster and cheaper. Also has an insanely large context window (it can handle a lot of data).",
			Icon:        "gemini",
			Name:        "gemini-2.0-flash",
			Provider:    llm.Gemini,
			Features: llm.ModelFeatures{
				HasVision:    true,
				HasPDF:       true,
				HasWebSearch: true,
				HasFast:      true,
			},
			Flags: llm.ModelFlags{
				IsFree: true,
			},
		})

		chatService.Handle(app.Router)

		if err = app.Start(); err != nil {
			fmt.Printf("Error starting the application: %v\n", err)
		}

	},
}
