package ollama

import (
	"errors"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/chat"
	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm/stream"
)

func StreamCompletion(req chat.Request) (*stream.Stream, error) {
	return nil, errors.New("not implemented")
}
