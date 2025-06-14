package ollama

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/stream"
	"github.com/ollama/ollama/api"
)

func ChatCompletion(req chat.Request) (*stream.Stream, error) {

	// TODO: replace with a global application client pool
	httpClient := &http.Client{
		Timeout: 0, // no global timeout; per-request ctx handles it
	}

	client := api.NewClient(
		&url.URL{
			Scheme: "https",
			Host:   "digitalfabrik.kilohertz.dev:3141",
			Path:   "/I5myvSQjxWZzQEKkn2z5sfrDvPLewIE54A58TAWUWUIrVZaC8WpV1bJENCgLQ4mhs9M15sstrzQfh80RjIw6Q9ezfaIIqQenvXaPh20SeNFiAzzmL9eLiQfPW3k5bx46",
		},
		httpClient,
	)

	doThink := true  // Defined by the user
	doStream := true // TODO: has to be true all the time for now

	ctx := context.TODO()
	request := &api.ChatRequest{
		Model:  "qwen3:30b",
		Think:  &doThink,
		Stream: &doStream,
		Options: map[string]any{
			"temperature": 0,
		},
		Messages: []api.Message{
			{
				Role:    "user",
				Content: req.Messages[0].Content,
			},
		},
	}

	s := stream.New()

	respFunc := func(resp api.ChatResponse) error {

		if resp.Done {
			s.Close()
		} else {
			s.Publish(stream.Chunk{
				Thinking: resp.Message.Thinking,
				Content:  resp.Message.Content,
			})
		}

		return nil

	}

	go func() {

		err := client.Chat(ctx, request, respFunc)
		if err != nil {
			s.Fail(err)
		}

	}()

	fmt.Println("Ollama completion started!") // TODO: Remove this debug statement

	return s, nil

}
