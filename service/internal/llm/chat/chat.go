package chat

// Request types
type Request struct {
	Model               string    `json:"model"`
	Temperature         float64   `json:"temperature,omitempty"`
	MaxCompletionTokens int       `json:"max_completion_tokens,omitempty"`
	TopP                float64   `json:"top_p,omitempty"`
	Stream              bool      `json:"stream"`
	ReasoningEffort     int32     `json:"reasoning_effort,omitempty"`
	Stop                any       `json:"stop,omitempty"`
	Messages            []Message `json:"messages"`
}

type Message struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Reasoning string `json:"reasoning"`
}
