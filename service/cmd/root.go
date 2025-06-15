package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version: fmt.Sprintf("%s, %s/%s", "0.0.1", runtime.GOOS, runtime.GOARCH),
	Use:     "t3-chat",
	Short:   "This application serves as a cloned rendition of the LLM Chat app originally developed by @theo, created in conjunction with the T3 Chat Cloneathon.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	// TODO: move this to the porper place
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

}
