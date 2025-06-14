package main

import (
	"github.com/Hinkolas/t3-chat-cloneathon/service/cmd"
	"github.com/joho/godotenv"
)

func main() {

	// TODO: move this to the porper place
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	cmd.Execute()
}
