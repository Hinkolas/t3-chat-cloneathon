package llm

import (
	"fmt"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/provider/anthropic"
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

func (mr *ModelRouter) ChatCompletion(req chat.Request) (*stream.Stream, error) {

	model, ok := mr.models[req.Model]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedModel, req.Model)
	}

	switch model.Provider {
	case Anthropic:
		// Handle Anthropic request
		return anthropic.ChatCompletion(req) // Assuming ollama is a package that
	case Ollama:
		// Handle Ollama request
		return ollama.ChatCompletion(req) // Assuming ollama is a package that
	default:
		// return "", fmt.Errorf("%w: %s", ErrUnsupportedProvider, model.Provider)
		panic(fmt.Errorf("%w: %s", ErrUnsupportedProvider, model.Provider))
	}

}
