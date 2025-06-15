package llm

import (
	"fmt"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/provider/anthropic"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/provider/gemini"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/provider/ollama"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/stream"
)

type ModelRouter struct {
	models map[string]Model
}

func NewModelRouter() *ModelRouter {
	return &ModelRouter{
		models: make(map[string]Model),
	}
}

func (mr *ModelRouter) AddModel(key string, model Model) {
	mr.models[key] = model
}

func (mr *ModelRouter) GetModel(key string) (Model, bool) {
	model, ok := mr.models[key]
	return model, ok
}

func (mr *ModelRouter) ListModels() map[string]Model {
	return mr.models
}

// TODO: Add Provider interface definition for the llm providers.
func (mr *ModelRouter) StreamCompletion(req chat.Request) (*stream.Stream, error) {

	// Get the model that was requested.
	// Return error if model does not exists.
	model, ok := mr.models[req.Model]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedModel, req.Model)
	}

	// Replace router with provider model name.
	req.Model = model.Name

	// Route the request to the corrosponding model provider.
	switch model.Provider {
	case Gemini:
		return gemini.StreamCompletion(req) // Handle request with Gemini
	case Anthropic:
		return anthropic.StreamCompletion(req) // Handle request with Ollama
	case Ollama:
		return ollama.StreamCompletion(req) // Handle request with Anthropic
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedModel, model.Provider) // THIS SHOULD NEVER HAPPEN!!!
	}

}
