package gemini

import (
	"context"
	"fmt"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/stream"
	"google.golang.org/genai"
)

func StreamCompletion(req chat.Request) (*stream.Stream, error) {

	ctx := context.TODO()

	// TODO: replace with a proper context
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*genai.Content, len(req.Messages)-1)
	// Convert universal format to gemini message format
	for i, message := range req.Messages[:len(req.Messages)-1] {
		messages[i] = &genai.Content{
			Parts: []*genai.Part{
				{Text: message.Content},
			},
		}
		if message.Role == "assistant" {
			messages[i].Role = genai.RoleModel
		} else if message.Role == "user" {
			messages[i].Role = genai.RoleUser
		}
	}

	chat, err := client.Chats.Create(ctx, req.Model, nil, messages)
	if err != nil {
		return nil, err
	}

	s := stream.New()

	go func() {

		for result, err := range chat.SendMessageStream(ctx, genai.Part{Text: req.Messages[len(req.Messages)-1].Content}) {

			s.Publish(stream.Chunk{
				Content: result.Text(),
			})

			if err != nil {
				fmt.Println("Gemini stream failed!") // TODO: Remove this debug statement
				s.Fail(err)
				return
			}

		}

		fmt.Println("Gemini stream completed!") // TODO: Remove this debug statement
		s.Close()

	}()

	fmt.Println("Gemini stream started!") // TODO: Remove this debug statement

	return s, nil

}
