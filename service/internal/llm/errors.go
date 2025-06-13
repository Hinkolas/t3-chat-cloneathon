package llm

import (
	"errors"
)

var (
	ErrUnsupportedProvider = errors.New("unsupported provider")
	ErrUnsupportedModel    = errors.New("unsupported model")
)
