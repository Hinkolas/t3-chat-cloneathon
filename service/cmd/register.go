package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/application/auth"
)

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP("config", "c", "config.yaml", "Path to config file")
	registerCmd.Flags().StringP("email", "e", "", "Email address")
	registerCmd.Flags().StringP("username", "u", "", "Username")
	registerCmd.Flags().StringP("password", "p", "", "Password")

	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Creates a new user with the given email, username and password",
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		// load config
		cfgFile, _ := cmd.Flags().GetString("config")
		cfg, err := application.LoadConfig(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		// bootstrap application
		app, err := application.NewApp(*cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing app: %v\n", err)
			os.Exit(1)
		}
		defer app.Close()

		// init auth service
		authService, err := auth.NewService(app)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing auth service: %v\n", err)
			os.Exit(1)
		}

		// create the user
		user, err := authService.CreateUser(context.Background(), email, username, password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating user: %v\n", err)
			os.Exit(1)
		}

		// print created user as JSON
		out, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting user JSON: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(out))
	},
}
