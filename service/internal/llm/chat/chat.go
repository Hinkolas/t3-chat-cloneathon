package chat

// type ChatProvider interface {
// 	ChatCompletion(ctx context.Context, req Request) (io.Reader, error)
// }

// Request types
type Request struct {
	Model               string    `json:"model"`
	Temperature         float64   `json:"temperature,omitempty"`
	MaxCompletionTokens int       `json:"max_completion_tokens,omitempty"`
	TopP                float64   `json:"top_p,omitempty"`
	Stream              bool      `json:"stream"`
	Thinking            int32     `json:"thinking,omitempty"`
	Stop                any       `json:"stop,omitempty"`
	Messages            []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response types
type Delta struct {
	Type    string `json:"type"` // "text" or "thinking"
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}
