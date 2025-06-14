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
	// TODO: implement chat completion logic using the Ollama SDK

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

	think := false

	ctx := context.Background()
	request := &api.ChatRequest{
		Model: "qwen3:30b",
		Think: &think,
		Messages: []api.Message{
			{
				Role:    "user",
				Content: req.Messages[0].Content,
			},
		},
	}

	s := stream.New()

	respFunc := func(resp api.ChatResponse) error {

		s.Publish(stream.Chunk{
			Thinking: resp.Message.Thinking,
			Content:  resp.Message.Content,
		})

		if resp.Done {
			s.Close()
			fmt.Println("Chat completion done!")
			// json.NewEncoder(os.Stdout).Encode(resp.Metrics)
		}
		return nil
	}

	go func() {

		err := client.Chat(ctx, request, respFunc)
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
			s.Close()
		}

	}()

	fmt.Println("Chat completion started!")

	return s, nil
}
