package anthropic

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/stream"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func StreamCompletion(req chat.Request, opt chat.Options) (*stream.Stream, error) {

	// TODO: handle invalid user provided api key or validate on input
	env := os.Getenv("ANTHROPIC_API_KEY")
	if user, ok := opt["anthropic_api_key"]; ok {
		env = user
	}

	if env == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is not set")
	}

	s := stream.New()

	// TODO: replace with a global application client pool
	httpClient := &http.Client{
		Timeout: 0, // no global timeout; per-request ctx handles it
	}

	client := anthropic.NewClient(
		option.WithAPIKey(env),
		option.WithHTTPClient(httpClient),
	)

	request := anthropic.MessageNewParams{
		Model:     anthropic.Model(req.Model),
		MaxTokens: int64(req.MaxCompletionTokens),
		Messages:  make([]anthropic.MessageParam, 0),
	}

	// Add temperature to the ollama request options
	if req.Temperature >= 0 && req.Temperature <= 1 {
		request.Temperature = anthropic.Float(req.Temperature)
	}

	if req.ReasoningEffort > 0 {
		request.Temperature = anthropic.Float(1)
		request.Thinking = anthropic.ThinkingConfigParamUnion{
			OfEnabled: &anthropic.ThinkingConfigEnabledParam{
				BudgetTokens: int64(req.ReasoningEffort),
				Type:         "enabled",
			},
		}
	}

	// Add system message to the ollama request messages
	if req.System != "" {
		request.System = []anthropic.TextBlockParam{
			{Text: req.System},
		}
	}

	for _, message := range req.Messages {
		blocks := []anthropic.ContentBlockParamUnion{
			anthropic.NewTextBlock(message.Content),
		}
		for _, attachment := range message.Attachments {
			if attachment.MimeType == "image/png" || attachment.MimeType == "image/jpeg" {
				blocks = append(blocks, anthropic.NewImageBlockBase64(attachment.MimeType, base64.StdEncoding.EncodeToString(attachment.Data)))
			} else if attachment.MimeType == "application/pdf" {
				blocks = append(blocks, anthropic.NewDocumentBlock(anthropic.Base64PDFSourceParam{
					Data:      base64.StdEncoding.EncodeToString(attachment.Data),
					MediaType: "application/pdf",
				}))
			}
		}
		if message.Role == "user" {
			request.Messages = append(request.Messages, anthropic.NewUserMessage(blocks...))
		} else if message.Role == "assistant" {
			request.Messages = append(request.Messages, anthropic.NewAssistantMessage(blocks...))
		}
	}

	go func() {

		completion := client.Messages.NewStreaming(s.Context(), request)

		message := anthropic.Message{}
		for completion.Next() {
			event := completion.Current()
			err := message.Accumulate(event)
			if err != nil {
				s.Fail(fmt.Errorf("anthropic: %w", err))
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
						Reasoning: deltaVariant.Thinking,
					})
				}
			}
		}

		if err := completion.Err(); err != nil {
			s.Fail(fmt.Errorf("anthropic: %w", err))
		} else {
			fmt.Println("Anthropic stream completed!") // TODO: Remove this debug statement
			s.Close()
		}

	}()

	fmt.Println("Anthropic stream started!") // TODO: Remove this debug statement

	return s, nil

}
