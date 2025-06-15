package anthropic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/stream"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func StreamCompletion(req chat.Request) (*stream.Stream, error) {

	// TODO: replace with a global application client pool
	httpClient := &http.Client{
		Timeout: 0, // no global timeout; per-request ctx handles it
	}

	client := anthropic.NewClient(
		option.WithHTTPClient(httpClient),
	)

	request := anthropic.MessageNewParams{
		Model:       anthropic.ModelClaude4Sonnet20250514,
		MaxTokens:   1024,
		Temperature: anthropic.Float(0.0),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(req.Messages[0].Content)),
		},
	}

	s := stream.New()

	go func() {

		completion := client.Messages.NewStreaming(context.TODO(), request)

		message := anthropic.Message{}
		for completion.Next() {
			event := completion.Current()
			err := message.Accumulate(event)
			if err != nil {
				s.Fail(err)
				return
			}

			switch eventVariant := event.AsAny().(type) {
			case anthropic.ContentBlockDeltaEvent:
				switch deltaVariant := eventVariant.Delta.AsAny().(type) {
				case anthropic.TextDelta:
					s.Publish(stream.Chunk{
						Content: deltaVariant.Text,
					})
				case anthropic.ThinkingDelta:
					s.Publish(stream.Chunk{
						Thinking: deltaVariant.Thinking,
					})
				}
			}
		}

		if err := completion.Err(); err != nil {
			s.Fail(err)
		} else {
			s.Close()
		}

	}()

	fmt.Println("Anthropic completion started!") // TODO: Remove this debug statement

	return s, nil

}
