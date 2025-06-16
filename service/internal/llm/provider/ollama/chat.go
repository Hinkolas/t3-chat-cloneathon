package ollama

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/stream"
	"github.com/ollama/ollama/api"
)

func StreamCompletion(req chat.Request, opt chat.Options) (*stream.Stream, error) {

	var base *url.URL
	var err error

	env, ok := os.LookupEnv("OLLAMA_BASE_URL")
	base, err = url.Parse(env)
	if !ok || err != nil {
		return nil, fmt.Errorf("failed to parse OLLAMA_BASE_URL: %w", err)
	}

	client := api.NewClient(base, &http.Client{}) // TODO: replace with a shared client pool

	request := &api.ChatRequest{
		Model:    req.Model,
		Think:    new(bool),
		Stream:   new(bool),
		Options:  make(map[string]any),
		Messages: make([]api.Message, len(req.Messages)),
	}

	*request.Think = req.ReasoningEffort > 0 // Defined by the user
	*request.Stream = true                   // TODO: has to be true all the time for now

	// Add temperature to the ollama request options
	if req.Temperature >= 0 && req.Temperature <= 1 {
		request.Options["temperature"] = req.Temperature
	}

	// Convert universal format to ollama message format
	for i, message := range req.Messages {
		request.Messages[i] = api.Message{
			Role:     message.Role,
			Content:  message.Content,
			Thinking: message.Reasoning, // TODO: Maybe remove to save ressources
		}
	}

	s := stream.New()

	respFunc := func(resp api.ChatResponse) error {

		// if the stream has been canceled, abort the Chat loop
		if err := s.Context().Err(); err != nil {
			return err
		}
		if resp.Done {
			s.Close()
		} else {
			s.Publish(stream.Chunk{
				Reasoning: resp.Message.Thinking,
				Content:   resp.Message.Content,
			})
		}
		return nil

	}

	go func() {

		err := client.Chat(s.Context(), request, respFunc) // TODO: replace with a proper context
		if err != nil {
			s.Fail(fmt.Errorf("ollama: %w", err))
		}

	}()

	fmt.Println("Ollama stream started!") // TODO: Remove this debug statement

	return s, nil

}
