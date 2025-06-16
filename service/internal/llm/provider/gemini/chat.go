package gemini

import (
	"context"
	"fmt"
	"log"
	"strings"

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

	config := genai.GenerateContentConfig{}

	// Add temperature to the ollama request options
	if req.Temperature >= 0 && req.Temperature <= 1 {
		temp := float32(req.Temperature)
		config.Temperature = &temp
	}

	if req.ReasoningEffort > 0 {
		config.ThinkingConfig = &genai.ThinkingConfig{
			IncludeThoughts: true,
			ThinkingBudget:  &req.ReasoningEffort,
		}
	}

	chat, err := client.Chats.Create(ctx, req.Model, &config, messages)
	if err != nil {
		return nil, err
	}

	s := stream.New()

	go func() {

		for result, err := range chat.SendMessageStream(ctx, genai.Part{Text: req.Messages[len(req.Messages)-1].Content}) {

			content, thoughts := getChunk(result)

			s.Publish(stream.Chunk{
				Content:   content,
				Reasoning: thoughts,
			})

			if err != nil {
				s.Fail(fmt.Errorf("gemini: %w", err))
				return
			}

		}

		fmt.Println("Gemini stream completed!") // TODO: Remove this debug statement
		s.Close()

	}()

	fmt.Println("Gemini stream started!") // TODO: Remove this debug statement

	return s, nil

}

func getChunk(r *genai.GenerateContentResponse) (string, string) {

	if len(r.Candidates) == 0 || r.Candidates[0].Content == nil || len(r.Candidates[0].Content.Parts) == 0 {
		return "", ""
	}

	if len(r.Candidates) > 1 {
		log.Println("Warning: there are multiple candidates in the response, returning text from the first one.")
	}

	var texts []string
	var thoughts []string
	var notTextParts []string
	for _, part := range r.Candidates[0].Content.Parts {
		if part.Text != "" {
			if part.Thought {
				thoughts = append(texts, part.Text)
			} else {
				texts = append(texts, part.Text)
			}
		} else {
			if part.InlineData != nil {
				notTextParts = append(notTextParts, "InlineData")
			}
			if part.CodeExecutionResult != nil {
				notTextParts = append(notTextParts, "CodeExecutionResult")
			}
			if part.ExecutableCode != nil {
				notTextParts = append(notTextParts, "ExecutableCode")
			}
			if part.FileData != nil {
				notTextParts = append(notTextParts, "FileData")
			}
			if part.FunctionCall != nil {
				notTextParts = append(notTextParts, "FunctionCall")
			}
			if part.FunctionResponse != nil {
				notTextParts = append(notTextParts, "FunctionResponse")
			}
		}
	}

	if len(notTextParts) > 0 {
		log.Printf("Warning: there are non-text parts %s in the response, returning concatenation of all text parts. Please refer to the non text parts for a full response from model.\n", strings.Join(notTextParts, ", "))
	}

	return strings.Join(texts, ""), strings.Join(thoughts, "")

}
