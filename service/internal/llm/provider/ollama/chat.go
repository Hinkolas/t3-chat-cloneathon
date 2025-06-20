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

	// TODO: handle invalid user provided api key or validate on input
	env := os.Getenv("OLLAMA_BASE_URL")
	if user, ok := opt["ollama_base_url"]; ok {
		env = user
	}

	if env == "" {
		return nil, fmt.Errorf("OLLAMA_BASE_URL is not set")
	}

	baseUrl, err := url.Parse(env)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OLLAMA_BASE_URL: %w", err)
	}

	client := api.NewClient(baseUrl, &http.Client{}) // TODO: replace with a shared client pool

	request := &api.ChatRequest{
		Model:    req.Model,
		Think:    new(bool),
		Stream:   new(bool),
		Options:  make(map[string]any),
		Messages: make([]api.Message, 0),
	}

	*request.Think = req.ReasoningEffort > 0 // Defined by the user
	*request.Stream = true                   // TODO: has to be true all the time for now

	// Add temperature to the ollama request options
	if req.Temperature >= 0 && req.Temperature <= 1 {
		request.Options["temperature"] = req.Temperature
	}

	// Add system message to the ollama request messages
	if req.System != "" {
		request.Messages = append(request.Messages, api.Message{
			Role:    "system",
			Content: req.System,
		})
	}

	// Convert universal format to ollama message format
	for _, message := range req.Messages {
		images := make([]api.ImageData, 0)
		for _, attachment := range message.Attachments {
			images = append(images, attachment.Data)
		}
		request.Messages = append(request.Messages, api.Message{
			Role:     message.Role,
			Content:  message.Content,
			Thinking: message.Reasoning, // TODO: Maybe remove to save ressources
			Images:   images,
		})
	}

	s := stream.New()

	respFunc := func(resp api.ChatResponse) error {

		// if the stream has been canceled, abort the Chat loop
		if err := s.Context().Err(); err != nil {
			return err
		}
		if resp.Done {
			fmt.Println("Ollama stream completed!")
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
